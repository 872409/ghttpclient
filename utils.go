package x_http_client

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

var sys_name string
var sys_release string
var sys_machine string

func init() {
	sys_name = runtime.GOOS
	sys_release = "-"
	sys_machine = runtime.GOARCH

	if out, err := exec.Command("uname", "-s").CombinedOutput(); err == nil {
		sys_name = string(bytes.TrimSpace(out))
	}
	if out, err := exec.Command("uname", "-r").CombinedOutput(); err == nil {
		sys_release = string(bytes.TrimSpace(out))
	}
	if out, err := exec.Command("uname", "-m").CombinedOutput(); err == nil {
		sys_machine = string(bytes.TrimSpace(out))
	}
}

func newTransport(conn *Conn, config *Config) *http.Transport {
	httpTimeOut := conn.config.HTTPTimeout
	httpMaxConns := conn.config.HTTPMaxConns
	// New Transport
	transport := &http.Transport{
		DialContext: func(ctx context.Context, netw, addr string) (net.Conn, error) {
			d := net.Dialer{
				Timeout:   httpTimeOut.ConnectTimeout,
				KeepAlive: 30 * time.Second,
			}
			if config.LocalAddr != nil {
				d.LocalAddr = config.LocalAddr
			}
			conn, err := d.DialContext(ctx, netw, addr)
			if err != nil {
				return nil, err
			}
			return newTimeoutConn(conn, httpTimeOut.ReadWriteTimeout, httpTimeOut.LongTimeout), nil
		},
		MaxIdleConns:          httpMaxConns.MaxIdleConns,
		MaxIdleConnsPerHost:   httpMaxConns.MaxIdleConnsPerHost,
		IdleConnTimeout:       httpTimeOut.IdleConnTimeout,
		ResponseHeaderTimeout: httpTimeOut.HeaderTimeout,
	}
	return transport
}



func calcMD5(body io.Reader, contentLen, md5Threshold int64) (reader io.Reader, b64 string, tempFile *os.File, err error) {
	if contentLen == 0 || contentLen > md5Threshold {
		// Huge body, use temporary file
		tempFile, err = ioutil.TempFile(os.TempDir(), TempFilePrefix)
		if tempFile != nil {
			io.Copy(tempFile, body)
			tempFile.Seek(0, os.SEEK_SET)
			md5 := md5.New()
			io.Copy(md5, tempFile)
			sum := md5.Sum(nil)
			b64 = base64.StdEncoding.EncodeToString(sum[:])
			tempFile.Seek(0, os.SEEK_SET)
			reader = tempFile
		}
	} else {
		// Small body, use memory
		buf, _ := ioutil.ReadAll(body)
		sum := md5.Sum(buf)
		b64 = base64.StdEncoding.EncodeToString(sum[:])
		reader = bytes.NewReader(buf)
	}
	return
}


// userAgent gets user agent
// It has the SDK version information, OS information and GO version
func userAgent() string {
	sys := getSysInfo()
	return fmt.Sprintf("wallet-trade-sdk-go/%s (%s/%s/%s;%s)", Version, sys.name,
		sys.release, sys.machine, runtime.Version())
}

type sysInfo struct {
	name    string // OS name such as windows/Linux
	release string // OS version 2.6.32-220.23.2.ali1089.el5.x86_64 etc
	machine string // CPU type amd64/x86_64
}

// getSysInfo gets system info
// gets the OS information and CPU type
func getSysInfo() sysInfo {
	return sysInfo{name: sys_name, release: sys_release, machine: sys_machine}
}




func GetReaderLen(reader io.Reader) (int64, error) {
	var contentLength int64
	var err error
	switch v := reader.(type) {
	case *bytes.Buffer:
		contentLength = int64(v.Len())
	case *bytes.Reader:
		contentLength = int64(v.Len())
	case *strings.Reader:
		contentLength = int64(v.Len())
	case *os.File:
		fInfo, fError := v.Stat()
		if fError != nil {
			err = fmt.Errorf("can't get reader content length,%s", fError.Error())
		} else {
			contentLength = fInfo.Size()
		}
	case *io.LimitedReader:
		contentLength = int64(v.N)
	case *LimitedReadCloser:
		contentLength = int64(v.N)
	default:
		err = fmt.Errorf("can't get reader content length,unkown reader type")
	}
	return contentLength, err
}

func LimitReadCloser(r io.Reader, n int64) io.Reader {
	var lc LimitedReadCloser
	lc.R = r
	lc.N = n
	return &lc
}

// LimitedRC support Close()
type LimitedReadCloser struct {
	io.LimitedReader
}

func (lc *LimitedReadCloser) Close() error {
	if closer, ok := lc.R.(io.ReadCloser); ok {
		return closer.Close()
	}
	return nil
}