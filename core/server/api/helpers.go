package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"reflect"
)

func AnswerError(c *gin.Context, statusCode int, message string) {
	switch statusCode {
	case 404:
		message = "not found. " + message
	case 500:
		message = "internal server error. " + message
	}

	c.JSON(statusCode, map[string]string{
		"status":        "error",
		"status_code":   fmt.Sprint(statusCode),
		"error_message": message,
	})
}

func AnswerResponse(c *gin.Context, data interface{}) {
	if reflect.TypeOf(data).Kind() == reflect.Ptr {
		data = reflect.ValueOf(data).Elem().Interface()
	}

	resp := map[string]interface{}{
		"status":      "ok",
		"status_code": http.StatusOK,
	}
	switch reflect.TypeOf(data).Kind() {
	case reflect.Slice, reflect.Array:
		if reflect.ValueOf(data).Len() == 0 {
			resp["results"] = []interface{}{}
		} else {
			resp["results"] = data
		}
		resp["number_of_results"] = reflect.ValueOf(data).Len()
	default:
		resp["result"] = data
	}
	c.JSON(http.StatusOK, resp)
}
