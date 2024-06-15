package middleware

import (
	"net/http"

	"github.com/entain-test-task/controller"
	"github.com/entain-test-task/model"
	"github.com/entain-test-task/repository"
	"github.com/gorilla/mux"
)

func Router(repository *repository.Store) *mux.Router {
	router := mux.NewRouter()

	controllers := controller.NewControllers(repository)

	router.Use(validateSourceHeader)

	router.HandleFunc("/users", controllers.User.GetAllUsers).Methods("GET")
	router.HandleFunc("/users/{user_id}", controllers.User.GetUserByID).Methods("GET")
	router.HandleFunc("/transactions/{user_id}", controllers.Transaction.GetAllTransactionsByUserID).Methods("GET")
	router.HandleFunc("/process-record/{user_id}", controllers.Transaction.ProcessRecord).Methods("POST")

	return router
}

func validateSourceHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !model.IsArrayContainsValue(model.SourceTypes, r.Header.Get("Source-Type")) {
			controller.StatusBadRequest(w, "Source-Type header is not valid")
			return
		}

		next.ServeHTTP(w, r)
	})
}
