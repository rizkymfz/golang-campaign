package helper

import (
	"math/rand"
	"time"

	"github.com/go-playground/validator/v10"
)

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

func RandomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, n)
	src := rand.NewSource(time.Now().UnixNano())

	for i, cache, remain := n-1, src.Int63(), 10; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), 10
		}
		if idx := int(cache & 0x3F); idx < len(letters) {
			result[i] = letters[idx]
			i--
		}
		cache >>= 6
		remain--
	}

	return string(result)
}
