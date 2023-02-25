package middleware

import (
	"encoding/json"
	"net/http"

	"github.com/entain-test-task/model"
)

func StatusOK(w http.ResponseWriter, response interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(response)
}

func StatusInternalServerError(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)

	json.NewEncoder(w).Encode(model.Error{
		Message: "Internal Server Error",
	})
}

func StatusBadRequest(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)

	json.NewEncoder(w).Encode(model.Error{
		Message: message,
	})
}

func StatusUnprocessableEntity(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnprocessableEntity)

	json.NewEncoder(w).Encode(model.Error{
		Message: message,
	})
}

func StatusBadRequestWithErrors(w http.ResponseWriter, message string, validationErrors []error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)

	errorDetails := make([]string, 0)
	for _, validationError := range validationErrors {
		errorDetails = append(errorDetails, validationError.Error())
	}

	json.NewEncoder(w).Encode(model.Error{
		Message: message,
		Errors:  errorDetails,
	})
}
