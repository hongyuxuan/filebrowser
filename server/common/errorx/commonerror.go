package errorx

import (
	"fmt"
	"net/http"
)

type BrowserError struct {
	Code    int         `json:"code"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type HttpErrorResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func NewError(code int, message string, data interface{}) error {
	return &BrowserError{Code: code, Message: message, Data: data}
}

func NewDefaultError(message string, a ...any) error {
	return &BrowserError{Code: http.StatusInternalServerError, Message: fmt.Sprintf(message, a...)}
}

func (e *BrowserError) Error() string {
	return e.Message
}

func (e *BrowserError) GetData() *HttpErrorResponse {
	return &HttpErrorResponse{
		Code:    e.Code,
		Message: e.Message,
		Data:    e.Data,
	}
}
