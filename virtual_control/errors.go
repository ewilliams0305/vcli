package vc4

import "fmt"

// Stores an invalid HTTP response as a status code.
type HttpResponseError struct {
	StatusCode int
}

func (code HttpResponseError) Error() string {
	return fmt.Sprintf("INVALID RESPONSE CODE RETURNED FROM SERVER %d", code)
}

func newResponseError(code int) *HttpResponseError {
	return &HttpResponseError{StatusCode: code}
}
