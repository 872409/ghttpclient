package x_http_client

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"hash"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Conn struct {
	config *Config
	url    *urlMaker
	client *http.Client
}

func (conn *Conn) init(config *Config, urlMaker *urlMaker, client *http.Client) error {
	if client == nil {
		// New transport
		transport := newTransport(conn, config)

		// Proxy
		if conn.config.IsUseProxy {
			proxyURL, err := url.Parse(config.ProxyHost)
			if err != nil {
				return err
			}
			if config.IsAuthProxy {
				if config.ProxyPassword != "" {
					proxyURL.User = url.UserPassword(config.ProxyUser, config.ProxyPassword)
				} else {
					proxyURL.User = url.User(config.ProxyUser)
				}
			}
			transport.Proxy = http.ProxyURL(proxyURL)
		}
		client = &http.Client{Transport: transport}
		if !config.RedirectEnabled {
			client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			}
		}
	}

	conn.config = config
	conn.url = urlMaker
	conn.client = client

	return nil
}

// Do sends request and returns the response
func (conn Conn) Do(method, path string, params map[string]interface{}, headers map[string]string, data io.Reader, listener ProgressListener) (*Response, error) {
	urlParams := conn.getURLParams(params)
	uri := conn.url.getURL(path, urlParams)
	return conn.doRequest(method, uri, headers, data, listener)
}

// Do sends request and returns the response
func (conn Conn) DoJSONResponse(method, path string, params map[string]interface{}, headers map[string]string, data interface{}, responseJSON interface{}) (*JSONResponse, error) {
	urlParams := conn.getURLParams(params)
	uri := conn.url.getURL(path, urlParams)
	resp, respErr := conn.doRequest(method, uri, headers, interface2JSONReader(data), nil)
	jsonResponse := &JSONResponse{Response: resp}

	jsonText, jsonErr := conn.jsonUnmarshal(resp.Body, responseJSON)
	jsonResponse.BodyJSONError = jsonErr
	jsonResponse.bodyText = jsonText

	return jsonResponse, respErr
}

func interface2JSONReader(v interface{}) *strings.Reader {
	bs, e := json.Marshal(v)
	if e != nil {
		return nil
	}
	// fmt.Printf("interface2JSONReader:%v \n", string(bs))
	return strings.NewReader(string(bs))
}
func (conn Conn) getURLParams(p interface{}) string {
	params, ok := p.(map[string]interface{})
	if !ok {
		return ""
	}

	// Sort
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Serialize
	var buf bytes.Buffer
	for _, k := range keys {
		if buf.Len() > 0 {
			buf.WriteByte('&')
		}
		buf.WriteString(url.QueryEscape(k))
		if params[k] != nil {
			buf.WriteString("=" + strings.Replace(url.QueryEscape(params[k].(string)), "+", "%20", -1))
		}
	}

	return buf.String()
}

func (conn Conn) doRequest(method string, uri *url.URL, headers map[string]string, data io.Reader, listener ProgressListener) (*Response, error) {
	method = strings.ToUpper(method)
	req := &http.Request{
		Method:     method,
		URL:        uri,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header:     make(http.Header),
		Host:       uri.Host,
	}

	tracker := &readerTracker{completedBytes: 0}
	fd, crc := conn.handleBody(req, data, listener, tracker)
	if fd != nil {
		defer func() {
			fd.Close()
			os.Remove(fd.Name())
		}()
	}

	if conn.config.IsAuthProxy {
		auth := conn.config.ProxyUser + ":" + conn.config.ProxyPassword
		basic := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
		req.Header.Set("Proxy-Authorization", basic)
	}

	date := time.Now().UTC().Format(http.TimeFormat)
	req.Header.Set(HTTPHeaderDate, date)
	req.Header.Set(HTTPHeaderHost, req.Host)
	req.Header.Set(HTTPHeaderUserAgent, conn.config.UserAgent)

	akIf := conn.config.GetCredentials()
	if akIf.GetSecurityToken() != "" {
		req.Header.Set(HTTPHeaderOssSecurityToken, akIf.GetSecurityToken())
	}

	if headers != nil {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
	}

	conn.signHeader(req)

	// Transfer started
	event := newProgressEvent(TransferStartedEvent, 0, req.ContentLength, 0)
	publishProgress(listener, event)

	if conn.config.LogLevel >= Debug {
		conn.LoggerHTTPReq(req)
	}

	resp, err := conn.client.Do(req)

	if err != nil {
		// Transfer failed
		event = newProgressEvent(TransferFailedEvent, tracker.completedBytes, req.ContentLength, 0)
		publishProgress(listener, event)
		conn.config.WriteLog(Debug, "[Resp:%p]http error:%s\n", req, err.Error())
		return nil, err
	}

	if conn.config.LogLevel >= Debug {
		// print out http resp
		conn.LoggerHTTPResp(req, resp)
	}

	// Transfer completed
	event = newProgressEvent(TransferCompletedEvent, tracker.completedBytes, req.ContentLength, 0)
	publishProgress(listener, event)

	return conn.handleResponse(resp, crc)
}

// handleResponse handles response
func (conn Conn) handleResponse(resp *http.Response, crc hash.Hash64) (*Response, error) {

	statusCode := resp.StatusCode
	if statusCode >= 400 && statusCode <= 505 {
		// 4xx and 5xx indicate that the operation has error occurred
		var respBody []byte
		respBody, err := readResponseBody(resp)
		if err != nil {
			return nil, err
		}

		if len(respBody) == 0 {
			err = ServiceError{
				StatusCode: statusCode,
				RequestID:  resp.Header.Get(HTTPHeaderOssRequestID),
			}
		} else {
			// Response contains storage service error object, unmarshal
			srvErr, errIn := serviceErrFromJSON(respBody, resp.StatusCode,
				resp.Header.Get(HTTPHeaderOssRequestID))
			if errIn != nil { // error unmarshaling the error response
				err = fmt.Errorf("service returned invalid response body, status = %s, RequestId = %s", resp.Status, resp.Header.Get(HTTPHeaderOssRequestID))
			} else {
				err = srvErr
			}
		}

		return &Response{
			StatusCode: resp.StatusCode,
			Headers:    resp.Header,
			Body:       ioutil.NopCloser(bytes.NewReader(respBody)), // restore the body
		}, err
	} else if statusCode >= 300 && statusCode <= 307 {
		// OSS use 3xx, but response has no body
		err := fmt.Errorf("service returned %d,%s", resp.StatusCode, resp.Status)
		return &Response{
			StatusCode: resp.StatusCode,
			Headers:    resp.Header,
			Body:       resp.Body,
		}, err
	}

	// 2xx, successful
	return &Response{
		StatusCode: resp.StatusCode,
		Headers:    resp.Header,
		Body:       resp.Body,
	}, nil
}

// LoggerHTTPReq Print the header information of the http request
func (conn Conn) LoggerHTTPReq(req *http.Request) {
	var logBuffer bytes.Buffer
	logBuffer.WriteString(fmt.Sprintf("[Req:%p]Method:%s\t", req, req.Method))
	logBuffer.WriteString(fmt.Sprintf("Host:%s\t", req.URL.Host))
	logBuffer.WriteString(fmt.Sprintf("Path:%s\t", req.URL.Path))
	logBuffer.WriteString(fmt.Sprintf("Query:%s\t", req.URL.RawQuery))
	logBuffer.WriteString(fmt.Sprintf("Header info:"))

	for k, v := range req.Header {
		var valueBuffer bytes.Buffer
		for j := 0; j < len(v); j++ {
			if j > 0 {
				valueBuffer.WriteString(" ")
			}
			valueBuffer.WriteString(v[j])
		}
		logBuffer.WriteString(fmt.Sprintf("\t%s:%s", k, valueBuffer.String()))
	}
	conn.config.WriteLog(Debug, "%s\n", logBuffer.String())
}

// LoggerHTTPResp Print Response to http request
func (conn Conn) LoggerHTTPResp(req *http.Request, resp *http.Response) {
	var logBuffer bytes.Buffer
	logBuffer.WriteString(fmt.Sprintf("[Resp:%p]StatusCode:%d\t", req, resp.StatusCode))
	logBuffer.WriteString(fmt.Sprintf("Header info:"))
	for k, v := range resp.Header {
		var valueBuffer bytes.Buffer
		for j := 0; j < len(v); j++ {
			if j > 0 {
				valueBuffer.WriteString(" ")
			}
			valueBuffer.WriteString(v[j])
		}
		logBuffer.WriteString(fmt.Sprintf("\t%s:%s", k, valueBuffer.String()))
	}
	conn.config.WriteLog(Debug, "%s\n", logBuffer.String())
}

// handleBody handles request body
func (conn Conn) handleBody(req *http.Request, body io.Reader, listener ProgressListener, tracker *readerTracker) (*os.File, hash.Hash64) {
	var file *os.File
	var crc hash.Hash64
	reader := body
	readerLen, err := GetReaderLen(reader)
	if err == nil {
		req.ContentLength = readerLen
	}
	req.Header.Set(HTTPHeaderContentLength, strconv.FormatInt(req.ContentLength, 10))

	// MD5
	if body != nil && conn.config.IsEnableMD5 && req.Header.Get(HTTPHeaderContentMD5) == "" {
		md5 := ""
		reader, md5, file, _ = calcMD5(body, req.ContentLength, conn.config.MD5Threshold)
		req.Header.Set(HTTPHeaderContentMD5, md5)
	}
	//
	// // CRC
	// if reader != nil && conn.config.IsEnableCRC {
	// 	crc = NewCRC(CrcTable(), initCRC)
	// 	reader = TeeReader(reader, crc, req.ContentLength, listener, tracker)
	// }

	// HTTP body
	rc, ok := reader.(io.ReadCloser)
	if !ok && reader != nil {
		rc = ioutil.NopCloser(reader)
	}

	// if conn.isUploadLimitReq(req) {
	// 	limitReader := &LimitSpeedReader{
	// 		reader:     rc,
	// 		ossLimiter: conn.config.UploadLimiter,
	// 	}
	// 	req.Body = limitReader
	// } else {
	// 	req.Body = rc
	// }
	req.Body = rc
	return file, crc
}

func readResponseBody(resp *http.Response) ([]byte, error) {
	defer resp.Body.Close()
	out, err := ioutil.ReadAll(resp.Body)
	if err == io.EOF {
		err = nil
	}
	return out, err
}

func serviceErrFromJSON(body []byte, statusCode int, requestID string) (ServiceError, error) {
	var storageErr ServiceError

	if err := json.Unmarshal(body, &storageErr); err != nil {
		return storageErr, err
	}

	storageErr.StatusCode = statusCode
	storageErr.RequestID = requestID
	storageErr.RawMessage = string(body)
	return storageErr, nil
}

func (conn Conn) jsonUnmarshal(body io.Reader, v interface{}) (string, error) {
	data, err := ioutil.ReadAll(body)
	if err != nil {
		return "", err
	}
	jsonStr := string(data)
	err = json.Unmarshal(data, v)
	if conn.config.LogLevel >= Debug {
		conn.config.WriteLog(Debug, "jsonUnmarshal:%v err:%v \n", jsonStr, err)
	}

	return jsonStr, err
}
