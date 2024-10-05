package http

import "fmt"

type HTTPError struct {
	Code    int
	Message string
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("Error %d: %s", e.Code, e.Message)
}

func NewHTTPError(code int, message string) *HTTPError {
	return &HTTPError{
		Code:    code,
		Message: message,
	}
}

var (
	errWrongPaginationQuery = NewHTTPError(120001, "Wrong pagination query")
	errWrongQuery           = NewHTTPError(120002, "Wrong query")
	errWrongBody            = NewHTTPError(120003, "Wrong body")
)
