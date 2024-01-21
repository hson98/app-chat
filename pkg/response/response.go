package response

import (
	"fmt"
	"strings"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Errors  interface{} `json:"errors"`
	Data    interface{} `json:"data"`
}

func BuildResponse(success bool, message string, data interface{}) Response {
	res := Response{
		Success: success,
		Message: message,
		Errors:  nil,
		Data:    data,
	}
	return res
}

func BuildErrorResponse(message string, err interface{}, data interface{}) Response {
	var splittedError interface{}
	switch err.(type) {
	case string:
		str := fmt.Sprintf("%v", err)
		splittedError = strings.Split(str, "\n")
	default:
		splittedError = err
	}

	res := Response{
		Success: false,
		Message: message,
		Errors:  splittedError,
		Data:    data,
	}
	return res
}
