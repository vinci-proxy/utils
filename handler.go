package utils

import (
	"io"
	"net"
	"net/http"
)

// ErrorHandler represents the error-specific interface required by error handlers.
type ErrorHandler interface {
	ServeHTTP(w http.ResponseWriter, req *http.Request, err error)
}

// StdHandler is an empty struct.
type StdHandler struct{}

// ErrorHandlerFunc represents function interface for error handlers.
type ErrorHandlerFunc func(http.ResponseWriter, *http.Request, error)

// DefaultHandler stores the default error handled to be used, which is an no-op.
var DefaultHandler ErrorHandler = &StdHandler{}

// ServeHTTP replies with the proper status code based on the given error and writes the body.
func (e *StdHandler) ServeHTTP(w http.ResponseWriter, req *http.Request, err error) {
	statusCode := http.StatusInternalServerError
	if e, ok := err.(net.Error); ok {
		if e.Timeout() {
			statusCode = http.StatusGatewayTimeout
		} else {
			statusCode = http.StatusBadGateway
		}
	} else if err == io.EOF {
		statusCode = http.StatusBadGateway
	}
	w.WriteHeader(statusCode)
	w.Write([]byte(http.StatusText(statusCode)))
}

// ServeHTTP calls f(w, r).
func (f ErrorHandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request, err error) {
	f(w, r, err)
}
