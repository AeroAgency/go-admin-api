package errors

import (
	"fmt"
	"net/http"
)

// Ошибка создания ресурса (уже был создан ранее)
type ConflictError struct {
	Err   error
	Trace string
}

func (e ConflictError) GetAppErrorCode() string {
	return "CONFLICT_ERROR"
}
func (e ConflictError) GetErrorLevel() string {
	return "error"
}
func (e ConflictError) GetAppErrorStatus() int {
	return http.StatusConflict
}

func (e ConflictError) GetAppErrorMessage() string {
	return "Некорректный запрос, ресурс был создан ранее"
}

func (e ConflictError) GetErrorDebugInfo() string {
	return fmt.Sprintf(" Service error: %v", e.Error())
}
func (e ConflictError) Error() string {
	return e.Err.Error()
}
func (e ConflictError) GetError() error {
	return e.Err
}
