package utils

import (
	"net/http"
)

const (
	InvalidHTTPStatus = ErrorMessageErr("invalid HTTP status")
	ErrGeneric  = ErrorMessageErr("generic error")
)

type ErrorMessageErr string

func (e ErrorMessageErr) Error() string {
	return string(e)
}

type ResponseWriter interface {
	PureJSON(status int, response interface{})
}

func DefaultErrorMessage(c ResponseWriter, httpStatus int, details interface{}) error {
	message := http.StatusText(httpStatus)
	if message == "" {
		return InvalidHTTPStatus
	}
	resp := map[string]interface{}{
		"code":    httpStatus,
		"message": http.StatusText(httpStatus),
		"details": details,
	}
	c.PureJSON(httpStatus, resp)
	return nil
}
