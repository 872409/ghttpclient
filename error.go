package x_http_client

import (
	"fmt"
)

// ServiceError contains fields of the error response from Oss Service REST API.
type ServiceError struct {
	Code       string   `json:"Code"`      // The error code returned from OSS to the caller
	Message    string   `json:"Message"`   // The detail error message from OSS
	RequestID  string   `json:"RequestId"` // The UUID used to uniquely identify the request
	HostID     string   `json:"HostId"`    // The OSS server cluster's Id
	Endpoint   string   `json:"Endpoint"`
	RawMessage string   // The raw messages from OSS
	StatusCode int      // HTTP status code
}

// Error implements interface error
func (e ServiceError) Error() string {
	if e.Endpoint == "" {
		return fmt.Sprintf("oss: service returned error: StatusCode=%d, ErrorCode=%s, ErrorMessage=\"%s\", RequestId=%s",
			e.StatusCode, e.Code, e.Message, e.RequestID)
	}
	return fmt.Sprintf("oss: service returned error: StatusCode=%d, ErrorCode=%s, ErrorMessage=\"%s\", RequestId=%s, Endpoint=%s",
		e.StatusCode, e.Code, e.Message, e.RequestID, e.Endpoint)
}