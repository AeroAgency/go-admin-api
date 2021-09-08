package errors

import (
	"fmt"
	"net/http"
)

//Ресурс не найден
type NotFoundError struct {
	Err   error
	Trace string
}

func (e NotFoundError) GetAppErrorCode() string {
	return "NOT_FOUND"
}
func (e NotFoundError) GetErrorLevel() string {
	return "error"
}
func (e NotFoundError) GetAppErrorStatus() int {
	return http.StatusNotFound
}

func (e NotFoundError) GetAppErrorMessage() string {
	return "Ресурс не найден"
}

func (e NotFoundError) GetErrorDebugInfo() string {
	return fmt.Sprintf("Stores service error: %v", e.Error())
}
func (e NotFoundError) Error() string {
	return e.Err.Error()
}
func (e NotFoundError) GetError() error {
	return e.Err
}
