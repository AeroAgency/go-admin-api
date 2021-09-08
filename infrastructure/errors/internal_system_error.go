package errors

import (
	"fmt"
	"net/http"
)

//Внутренняя ошибка сервиса, определенная на уровне инфраструктуры
type InternalSystemError struct {
	Err     error
	Message string
	Trace   string
}

func (e InternalSystemError) GetErrorLevel() string {
	return "error"
}
func (e InternalSystemError) GetAppErrorCode() string {
	return "INTERNAL_SERVER_ERROR"
}

func (e InternalSystemError) GetAppErrorStatus() int {
	return http.StatusInternalServerError
}

func (e InternalSystemError) GetAppErrorMessage() string {
	return "Произошла ошибка. Попробуйте выполнить операцию позже"
}

func (e InternalSystemError) GetErrorDebugInfo() string {
	return fmt.Sprintf("Stores control service error: %v: %v", e.Message, e.Error())
}

func (e InternalSystemError) Error() string {
	return e.Err.Error()
}
func (e InternalSystemError) GetError() error {
	return e.Err
}
