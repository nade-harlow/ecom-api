package apperrors

import "net/http"

type AppError interface {
	Error() string
	GetCode() int
}

type httpError struct {
	message string
	code    int
}

func (err *httpError) Error() string {
	return err.message
}
func (err *httpError) GetCode() int {
	return err.code
}

func BadRequestError(message string) *httpError {
	return &httpError{
		message: message,
		code:    http.StatusBadRequest,
	}
}

func ConflictError(message string) *httpError {
	return &httpError{
		message: message,
		code:    http.StatusConflict,
	}
}

func ForbiddenError(message string) *httpError {
	return &httpError{
		message: message,
		code:    http.StatusForbidden,
	}
}

func InternalServerError(message string) *httpError {
	return &httpError{
		message: message,
		code:    http.StatusInternalServerError,
	}
}

func NotFoundError(message string) *httpError {
	return &httpError{
		message: message,
		code:    http.StatusNotFound,
	}
}

func UnauthorizedError(message string) *httpError {
	return &httpError{
		message: message,
		code:    http.StatusUnauthorized,
	}
}
