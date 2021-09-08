package errors

import (
	"fmt"
	"net/http"
)

type BadRequestError struct {
	Err   error
	Trace string
}

func (e BadRequestError) GetAppErrorCode() string {
	return "BAD_REQUEST"
}
func (e BadRequestError) GetErrorLevel() string {
	return "error"
}
func (e BadRequestError) GetAppErrorStatus() int {
	return http.StatusBadRequest
}

func (e BadRequestError) GetAppErrorMessage() string {
	return "Сервер не может распознать запрос"
}

func (e BadRequestError) GetErrorDebugInfo() string {
	return fmt.Sprintf("Stores service error: %v", e.Error())
}
func (e BadRequestError) Error() string {
	return e.Err.Error()
}
func (e BadRequestError) GetError() error {
	return e.Err
}
