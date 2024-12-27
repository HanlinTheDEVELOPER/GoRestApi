package main

import (
	"encoding/json"
	"net/http"
)

func SuccessResponse(writer http.ResponseWriter, statusCode int, payload interface{}) {
	response, _ := json.Marshal(payload)
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(statusCode)
	_, _ = writer.Write(response)
}

func ErrorResponse(writer http.ResponseWriter, statusCode int, error string) {
	errorMessage := map[string]string{"error": error}
	SuccessResponse(writer, statusCode, errorMessage)
}
