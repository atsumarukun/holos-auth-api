package apierr

type ApiError interface {
	Error() (int, string)
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

func (e *apiError) Error() (int, string) {
	return e.code, e.message
}

func (e *apiError) Code() int {
	return e.code
}

func (e apiError) Message() string {
	return e.message
}
