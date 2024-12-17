package apierr

import "fmt"

type ApiError interface {
	Error() string
	Code() int
	Message() string
}

type apiError struct {
	code    int
	message string
}

func NewApiError(code int, message string) ApiError {
	return &apiError{
		code:    code,
		message: message,
	}
}

func (e *apiError) Error() string {
	return fmt.Sprintf("code: %d, message: %s", e.code, e.message)
}

func (e *apiError) Code() int {
	return e.code
}

func (e apiError) Message() string {
	return e.message
}
