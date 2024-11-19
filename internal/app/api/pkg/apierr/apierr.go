package apierr

type ApiError interface {
	Error() (int, string)
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

func (ae *apiError) Error() (int, string) {
	return ae.code, ae.message
}
