package x_http_client

import (
	"fmt"
)

// ServiceError contains fields of the error response from Oss Service REST API.
type ServiceError struct {
	Code       int    `json:"code"`       // The error code returned from server to the caller
	Message    string `json:"msg"`        // The detail error message from server
	TrackID    string `json:"track_id"`   // The UUID used to uniquely identify the request
	RequestID  string `json:"request_id"` // The UUID used to uniquely identify the request
	HostID     string `json:"HostId"`     // The  server cluster's Id
	Endpoint   string `json:"Endpoint"`
	RawMessage string // The raw messages
	StatusCode int    // HTTP status code
}

// Error implements interface error
func (e ServiceError) Error() string {
	if e.Endpoint == "" {
		return fmt.Sprintf("service returned error: StatusCode=%d, ErrorCode=%s, ErrorMessage=\"%s\", RequestId=%s",
			e.StatusCode, e.Code, e.Message, e.RequestID)
	}
	return fmt.Sprintf("service returned error: StatusCode=%d, ErrorCode=%s, ErrorMessage=\"%s\", RequestId=%s, Endpoint=%s",
		e.StatusCode, e.Code, e.Message, e.RequestID, e.Endpoint)
}
