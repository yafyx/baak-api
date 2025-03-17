package utils

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func WriteJSONResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{
		Success: true,
		Data:    data,
	})
}

func WriteErrorResponse(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(Response{
		Success: false,
		Error:   message,
	})
}

func WriteValidationError(w http.ResponseWriter, message string) {
	WriteErrorResponse(w, http.StatusBadRequest, message)
}

func WriteNotFoundError(w http.ResponseWriter) {
	WriteErrorResponse(w, http.StatusNotFound, "Resource not found")
}

func WriteInternalServerError(w http.ResponseWriter) {
	WriteErrorResponse(w, http.StatusInternalServerError, "Internal server error")
}
