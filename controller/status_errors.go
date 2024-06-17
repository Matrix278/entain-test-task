package controller

import (
	"encoding/json"
	"net/http"

	"github.com/entain-test-task/model"
)

func StatusOK(w http.ResponseWriter, response interface{}) {
	writeJSONResponse(w, http.StatusOK, response)
}

func StatusInternalServerError(w http.ResponseWriter) {
	writeJSONResponse(w, http.StatusInternalServerError, model.Error{
		Message: "Internal Server Error",
	})
}

func StatusBadRequest(w http.ResponseWriter, message string) {
	writeJSONResponse(w, http.StatusBadRequest, model.Error{
		Message: message,
	})
}

func StatusUnprocessableEntity(w http.ResponseWriter, message string) {
	writeJSONResponse(w, http.StatusUnprocessableEntity, model.Error{
		Message: message,
	})
}

func StatusBadRequestWithErrors(w http.ResponseWriter, message string, validationErrors []error) {
	errorDetails := make([]string, 0)
	for _, validationError := range validationErrors {
		errorDetails = append(errorDetails, validationError.Error())
	}

	writeJSONResponse(w, http.StatusBadRequest, model.Error{
		Message: message,
		Errors:  errorDetails,
	})
}

func writeJSONResponse(w http.ResponseWriter, statusCode int, response interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
