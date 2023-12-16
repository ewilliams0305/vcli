package vc

import (
	"fmt"
)

type VirtualControlError interface {
	Error() string
}

type ServerError struct {
	code int
	err  error
}

func (e *ServerError) Error() string {
	return fmt.Sprintf("%d | %s", e.code, e.err)
}

func NewServerError(code int, err error) *ServerError {
	return &ServerError{
		err:  err,
		code: code,
	}
}

// Stores an invalid HTTP response as a status code.
type HttpResponseError struct {
	StatusCode int
}

func (code HttpResponseError) Error() string {
	return fmt.Sprintf("INVALID RESPONSE CODE RETURNED FROM SERVER %d", code)
}

func NewResponseError(code int) *HttpResponseError {
	return &HttpResponseError{StatusCode: code}
}
