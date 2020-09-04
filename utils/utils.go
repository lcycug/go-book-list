package utils

import (
	"book-list/models"
	"encoding/json"
	"net/http"
)

func SendError(w http.ResponseWriter, status int, err models.Error) {
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(err)
}

func SendSuccess(w http.ResponseWriter, data interface{}) {
	_ = json.NewEncoder(w).Encode(data)
}
