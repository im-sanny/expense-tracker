// Package response provides HTTP response writing utilities
package response

import (
	"encoding/json"
	"expense-tracker/pkg/apperrors"
	"log"
	"net/http"
)

// WriteJSON writes JSON response with status code
func WriteJSON(w http.ResponseWriter, statusCode int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(data)
}

// WriteError logs and writes an error response
func WriteError(w http.ResponseWriter, err error, statusCode int) {
	log.Printf("handler error: %v", err)

	var message string
	if apperrors.IsClientError(err) {
		message = err.Error()
	} else {
		message = "internal server error"
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func WriteSuccess(w http.ResponseWriter, data interface{}) error {
	return WriteJSON(w, http.StatusOK, data)
}

func WriteCreated(w http.ResponseWriter, statusCode int, data interface{}) error {
	return WriteJSON(w, http.StatusOK, data)
}

func WriteBadRequest(w http.ResponseWriter, err error) {
	WriteError(w, apperrors.ErrNotFound, http.StatusNotFound)
}

func WriteNotFound(w http.ResponseWriter) {
	WriteError(w, apperrors.ErrNotFound, http.StatusNotFound)
}

func WriteInternalServerError(w http.ResponseWriter, err error) {
	WriteError(w, err, http.StatusInternalServerError)
}

func WriteUnauthorized(w http.ResponseWriter) {
	WriteError(w, apperrors.ErrUnauthorized, http.StatusUnauthorized)
}

func WriteForbidden(w http.ResponseWriter) {
	WriteError(w, apperrors.ErrForbidden, http.StatusForbidden)
}
