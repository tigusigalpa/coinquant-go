package coinquant

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// APIError is a typed error returned by the CoinQuant API.
type APIError struct {
	Status    int    `json:"-"`
	RequestID string `json:"request_id"`
	Code      string `json:"code"`
	Message   string `json:"message"`
}

// Error implements the error interface.
func (e *APIError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("coinquant API error status=%d request_id=%s code=%s message=%s", e.Status, e.RequestID, e.Code, e.Message)
	}
	return fmt.Sprintf("coinquant API error status=%d request_id=%s", e.Status, e.RequestID)
}

// newAPIError builds an APIError from an HTTP response body and status.
func newAPIError(status int, body []byte, requestID string) error {
	err := &APIError{Status: status, RequestID: requestID}
	_ = json.Unmarshal(body, err)
	// if the error body is wrapped, try to unwrap
	if err.Code == "" && err.Message == "" {
		var wrapped struct {
			Error *APIError `json:"error"`
		}
		if json.Unmarshal(body, &wrapped) == nil && wrapped.Error != nil {
			err.Code = wrapped.Error.Code
			err.Message = wrapped.Error.Message
			err.RequestID = wrapped.Error.RequestID
		}
	}
	if err.RequestID == "" {
		err.RequestID = requestID
	}
	return err
}

// IsAPIError reports whether err is or wraps a *APIError.
func IsAPIError(err error) (*APIError, bool) {
	for err != nil {
		if e, ok := err.(*APIError); ok {
			return e, true
		}
		if u, ok := err.(interface{ Unwrap() error }); ok {
			err = u.Unwrap()
			continue
		}
		break
	}
	return nil, false
}

// HTTPResponse wraps the raw *http.Response for consumers who need it.
type HTTPResponse struct {
	*http.Response
	RequestID string
}
