package errors

import (
	"fmt"
	"net/http"
)

//Внутренняя ошибка сервиса, определенная на уровне сервиса
type InternalError struct {
	Err   error
	Trace string
}

func (e InternalError) GetAppErrorCode() string {
	return "INTERNAL_SERVER_ERROR"
}
func (e InternalError) GetErrorLevel() string {
	return "error"
}
func (e InternalError) GetAppErrorStatus() int {
	return http.StatusInternalServerError
}

func (e InternalError) GetAppErrorMessage() string {
	return "Произошла ошибка. Попробуйте выполнить операцию позже"
}

func (e InternalError) GetErrorDebugInfo() string {
	return fmt.Sprintf("Stores service error: %v", e.Error())
}
func (e InternalError) Error() string {
	return e.Err.Error()
}
func (e InternalError) GetError() error {
	return e.Err
}
