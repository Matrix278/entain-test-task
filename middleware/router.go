package middleware

import (
	"github.com/gorilla/mux"
)

func Router(handler *Handler) *mux.Router {
	router := mux.NewRouter()

	router.Use(ValidateSourceHeader)

	router.HandleFunc("/users", handler.GetAllUsers).Methods("GET")
	router.HandleFunc("/users/{user_id}", handler.GetUserByID).Methods("GET")
	router.HandleFunc("/transactions/{user_id}", handler.GetAllTransactionsByUserID).Methods("GET")
	router.HandleFunc("/process-record/{user_id}", handler.ProcessRecord).Methods("POST")

	return router
}
