// from docker's httputils
package httputils

import (
	"encoding/json"
	"net/http"
)

// httpStatusError is an interface
// that errors with custom status codes
// implement to tell the api layer
// which response status to set.
type httpStatusError interface {
	HTTPErrorStatusCode() int
}

// inputValidationError is an interface
// that errors generated by invalid
// inputs can implement to tell the
// api layer to set a 400 status code
// in the response.
type inputValidationError interface {
	IsValidationError() bool
}

type errorResponse struct {
	Error string `json:"error"`
}

// WriteError decodes a specific docker error and sends it in the response.
func WriteError(w http.ResponseWriter, err error) {
	if err == nil || w == nil {
		return
	}

	var statusCode int
	errMsg := err.Error()

	switch e := err.(type) {
	case httpStatusError:
		statusCode = e.HTTPErrorStatusCode()
	case inputValidationError:
		statusCode = http.StatusBadRequest
	}

	if statusCode == 0 {
		statusCode = http.StatusInternalServerError
	}

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(errorResponse{errMsg})
}