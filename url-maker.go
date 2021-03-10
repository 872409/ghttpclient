package x_http_client

import (
	"fmt"
	"net"
	"net/url"
	"strings"
	"time"
)

const (
	urlTypeIP     = 2
	urlTypeAliyun = 3
)

type urlMaker struct {
	Scheme string // HTTP or HTTPS
	NetLoc string // Host or IP
	Type   int    //  2 IP, 3 ALIYUN
	// IsProxy bool   // Proxy
}

// getURL gets URL
func (um urlMaker) getURL(path, params string) *url.URL {
	host := um.NetLoc

	addr := ""
	if params == "" {
		addr = fmt.Sprintf("%s://%s%s", um.Scheme, host, path)
	} else {
		addr = fmt.Sprintf("%s://%s%s?%s", um.Scheme, host, path, params)
	}
	uri, _ := url.ParseRequestURI(addr)
	return uri
}

func (um *urlMaker) Init(endpoint string) error {
	if strings.HasPrefix(endpoint, "http://") {
		um.Scheme = "http"
		um.NetLoc = endpoint[len("http://"):]
	} else if strings.HasPrefix(endpoint, "https://") {
		um.Scheme = "https"
		um.NetLoc = endpoint[len("https://"):]
	} else {
		um.Scheme = "http"
		um.NetLoc = endpoint
	}

	// use url.Parse() to get real host
	strUrl := um.Scheme + "://" + um.NetLoc
	url, err := url.Parse(strUrl)
	if err != nil {
		return err
	}

	um.NetLoc = url.Host
	host, _, err := net.SplitHostPort(um.NetLoc)
	if err != nil {
		host = um.NetLoc
		if host[0] == '[' && host[len(host)-1] == ']' {
			host = host[1 : len(host)-1]
		}
	}

	ip := net.ParseIP(host)
	if ip != nil {
		um.Type = urlTypeIP
	} else {
		um.Type = urlTypeAliyun
	}
	// um.IsProxy = isProxy

	return nil
}

// timeoutConn handles HTTP timeout
type timeoutConn struct {
	conn        net.Conn
	timeout     time.Duration
	longTimeout time.Duration
}

func newTimeoutConn(conn net.Conn, timeout time.Duration, longTimeout time.Duration) *timeoutConn {
	conn.SetReadDeadline(time.Now().Add(longTimeout))
	return &timeoutConn{
		conn:        conn,
		timeout:     timeout,
		longTimeout: longTimeout,
	}
}

func (c *timeoutConn) Read(b []byte) (n int, err error) {
	c.SetReadDeadline(time.Now().Add(c.timeout))
	n, err = c.conn.Read(b)
	c.SetReadDeadline(time.Now().Add(c.longTimeout))
	return n, err
}

func (c *timeoutConn) Write(b []byte) (n int, err error) {
	c.SetWriteDeadline(time.Now().Add(c.timeout))
	n, err = c.conn.Write(b)
	c.SetReadDeadline(time.Now().Add(c.longTimeout))
	return n, err
}

func (c *timeoutConn) Close() error {
	return c.conn.Close()
}

func (c *timeoutConn) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}

func (c *timeoutConn) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

func (c *timeoutConn) SetDeadline(t time.Time) error {
	return c.conn.SetDeadline(t)
}

func (c *timeoutConn) SetReadDeadline(t time.Time) error {
	return c.conn.SetReadDeadline(t)
}

func (c *timeoutConn) SetWriteDeadline(t time.Time) error {
	return c.conn.SetWriteDeadline(t)
}

type readerTracker struct {
	completedBytes int64
}

// ProgressEventType defines transfer progress event type
type ProgressEventType int

const (
	// TransferStartedEvent transfer started, set TotalBytes
	TransferStartedEvent ProgressEventType = 1 + iota
	// TransferDataEvent transfer data, set ConsumedBytes anmd TotalBytes
	TransferDataEvent
	// TransferCompletedEvent transfer completed
	TransferCompletedEvent
	// TransferFailedEvent transfer encounters an error
	TransferFailedEvent
)

// ProgressEvent defines progress event
type ProgressEvent struct {
	ConsumedBytes int64
	TotalBytes    int64
	RwBytes       int64
	EventType     ProgressEventType
}

// ProgressListener listens progress change
type ProgressListener interface {
	ProgressChanged(event *ProgressEvent)
}

// -------------------- Private --------------------

func newProgressEvent(eventType ProgressEventType, consumed, total int64, rwBytes int64) *ProgressEvent {
	return &ProgressEvent{
		ConsumedBytes: consumed,
		TotalBytes:    total,
		RwBytes:       rwBytes,
		EventType:     eventType}
}

// publishProgress
func publishProgress(listener ProgressListener, event *ProgressEvent) {
	if listener != nil && event != nil {
		listener.ProgressChanged(event)
	}
}
