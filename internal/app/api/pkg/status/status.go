package status

import (
	"fmt"
	"net/http"
)

type Status struct {
	code    int
	message string
}

func New(code int, message string) *Status {
	return &Status{
		code:    code,
		message: message,
	}
}

func Error(code int, message string) error {
	return &Status{
		code:    code,
		message: message,
	}
}

func FromError(err error) *Status {
	if err == nil {
		return nil
	}

	if v, ok := err.(*Status); ok {
		return v
	} else {
		return &Status{
			code:    http.StatusInternalServerError,
			message: err.Error(),
		}
	}
}

func (e *Status) Error() string {
	return fmt.Sprintf("code: %d, message: %s", e.code, e.message)
}

func (e *Status) Code() int {
	return e.code
}

func (e *Status) Message() string {
	return e.message
}
