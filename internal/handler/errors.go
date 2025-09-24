package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/DiDinar5/1inch_test_task/domain"
)

func errrorJson(statusCode int, message string, writer http.ResponseWriter) {
	writer.WriteHeader(statusCode)
	resp, err := json.Marshal(&domain.CommonResponse{
		Message: message,
		Status:  false,
	})
	if err != nil {
		log.Println(err)
	}

	writer.Write(resp)
}
