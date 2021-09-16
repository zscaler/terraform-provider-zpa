package client

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

type ErrorResponse struct {
	Response *http.Response
	Message  string
}

func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("FAILED: %v, %v, %d, %v, %v", r.Response.Request.Method, r.Response.Request.URL, r.Response.StatusCode, r.Response.Status, r.Message)
}

func checkErrorInResponse(res *http.Response) error {
	if c := res.StatusCode; c >= 200 && c <= 299 {
		return nil
	}
	errorResponse := &ErrorResponse{Response: res}
	errorMessage, err := ioutil.ReadAll(res.Body)
	if err == nil && len(errorMessage) > 0 {
		errorResponse.Message = string(errorMessage)
	}
	return errorResponse
}

// IsObjectNotFound returns true on missing object error (404).
func (r ErrorResponse) IsObjectNotFound() bool {
	return r.Response.StatusCode == 404
}
