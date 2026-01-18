package cerrors

import "net/http"

const (
	notFoundCode       = "NOT_FOUND"
	badRequestCode     = "BAD_REQUEST"
	internalServerCode = "INTERNAL_SERVER"
)

type AppError struct {
	Status  int    `json:"-"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *AppError) Error() string {
	return e.Message
}

func NotFound(msg string) *AppError {
	return &AppError{
		Status:  http.StatusNotFound,
		Code:    notFoundCode,
		Message: msg,
	}
}

func BadRequest(msg string) *AppError {
	return &AppError{
		Status:  http.StatusBadRequest,
		Code:    badRequestCode,
		Message: msg,
	}
}

func Internal(msg string) *AppError {
	return &AppError{
		Status:  http.StatusInternalServerError,
		Code:    internalServerCode,
		Message: msg,
	}
}
