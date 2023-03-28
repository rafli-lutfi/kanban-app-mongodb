package models

import (
	"encoding/json"
	"net/http"
)

func ResponeWithError(w http.ResponseWriter, statusCode int, msg string) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": statusCode,
		"error":  msg,
	})
}

func ResponeWithJson(w http.ResponseWriter, statusCode int, msg string, payload interface{}) {
	w.WriteHeader(statusCode)

	if payload == nil {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  statusCode,
			"message": msg,
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  statusCode,
		"message": msg,
		"data":    payload,
	})
}
