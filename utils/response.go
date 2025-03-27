package utils

import (
	"encoding/json"
	"net/http"
	"strings"
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

func WriteHTTPError(w http.ResponseWriter, err error) {
	message := err.Error()

	// Handle specific HTTP errors with appropriate status codes
	if strings.Contains(message, "unexpected status code: 403") ||
		strings.Contains(message, "access forbidden (403)") {
		WriteErrorResponse(w, http.StatusServiceUnavailable,
			"Service temporarily unavailable due to access restrictions. Please try again later.")
		return
	} else if strings.Contains(message, "unexpected status code: 429") {
		WriteErrorResponse(w, http.StatusTooManyRequests,
			"Too many requests to the backend server. Please try again later.")
		return
	} else if strings.Contains(message, "unexpected status code: 5") {
		// Handle any 5xx error
		WriteErrorResponse(w, http.StatusBadGateway,
			"The backend server is experiencing issues. Please try again later.")
		return
	}

	// Default to internal server error for other cases
	WriteInternalServerError(w)
}
