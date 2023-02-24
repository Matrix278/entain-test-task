package middleware

import (
	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	router := mux.NewRouter()

	router.Use(ValidateSourceHeader)

	router.HandleFunc("/users", GetAllUsers).Methods("GET")
	router.HandleFunc("/users/{user_id}", GetUserByID).Methods("GET")
	router.HandleFunc("/transactions/{user_id}", GetAllTransactionsByUserID).Methods("GET")
	router.HandleFunc("/process-record/{user_id}", ProcessRecord).Methods("POST")

	return router
}
