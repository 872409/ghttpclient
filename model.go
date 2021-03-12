package x_http_client

import (
	"io"
	"io/ioutil"
	"net/http"
)

// Response defines HTTP response from OSS
type Response struct {
	StatusCode int
	RequestID  string
	TrackID    string

	Headers        http.Header
	Body           io.ReadCloser
	bodyText       string
	isBodyTextRead bool
	// ClientCRC  uint64
	// ServerCRC  uint64
}

func (r *Response) Read(p []byte) (n int, err error) {
	return r.Body.Read(p)
}

// Close close http reponse body
func (r *Response) Close() error {
	return r.Body.Close()
}

func (r *Response) GetBodyText() string {
	if r.isBodyTextRead || len(r.bodyText) > 0 {
		return r.bodyText
	}
	r.isBodyTextRead = true
	b, e := ioutil.ReadAll(r.Body)
	if e == nil {
		r.bodyText = string(b)
	}
	return r.bodyText
}

type JSONResponse struct {
	*Response
	BodyJSONError error
	// BodyJSON      interface{}
}
