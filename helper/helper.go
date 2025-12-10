package helper

import "github.com/go-playground/validator/v10"

type response struct {
	Status  string `json:"status"`
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func ApiResponse(status string, code int, message string, data any) response {
	jsonResponse := response{
		Status:  status,
		Code:    code,
		Message: message,
		Data:    data,
	}

	return jsonResponse
}

func FormatValidationError(err error) []string {
	var errors []string
	for _, e := range err.(validator.ValidationErrors) {
		errors = append(errors, e.Error())
	}

	return errors
}
