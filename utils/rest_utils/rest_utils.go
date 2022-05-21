package rest_utils

import (
	"net/http"
)

type RestResponseDTO struct {
	Data      interface{} `json:"data,omitempty"`
	Code      int         `json:"code,omitempty"`
	Message   string      `json:"message,omitempty"`
	Error     string      `json:"error,omitempty"`
	ErrorCode int         `json:"error_code,omitempty"`
}

func PrepareApiResponseAsJson(w http.ResponseWriter) error {
	w.Header().Add("Content-Type", "application/json")
	return nil
}

func NewApiResponse(code int, data interface{}, message string) RestResponseDTO {
	response := RestResponseDTO{
		Data:    data,
		Code:    code,
		Message: message,
	}
	return response
}

func NewErrorApiResponse(code int, errorCode int, error string) RestResponseDTO {
	response := RestResponseDTO{
		Code:      code,
		Error:     error,
		ErrorCode: errorCode,
	}

	return response
}
