package middleware

import (
	"net/http"

	"github.com/entain-test-task/model"
	"github.com/gorilla/mux"
)

func Router(handler *Handler) *mux.Router {
	router := mux.NewRouter()

	router.Use(validateSourceHeader)

	router.HandleFunc("/users", handler.GetAllUsers).Methods("GET")
	router.HandleFunc("/users/{user_id}", handler.GetUserByID).Methods("GET")
	router.HandleFunc("/transactions/{user_id}", handler.GetAllTransactionsByUserID).Methods("GET")
	router.HandleFunc("/process-record/{user_id}", handler.ProcessRecord).Methods("POST")

	return router
}

func validateSourceHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !model.IsArrayContainsValue(model.SourceTypes, r.Header.Get("Source-Type")) {
			StatusBadRequest(w, "Source-Type header is not valid")
			return
		}

		next.ServeHTTP(w, r)
	})
}
