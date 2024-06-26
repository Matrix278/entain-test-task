package middleware

import (
	"net/http"

	"github.com/entain-test-task/controller"
	"github.com/entain-test-task/model"
	"github.com/entain-test-task/model/enum"
	"github.com/gorilla/mux"
)

func Router(controllers *controller.Controllers) *mux.Router {
	router := mux.NewRouter()

	router.Use(validateSourceHeader)

	router.HandleFunc("/users", controllers.User.GetAllUsers).Methods(http.MethodGet)
	router.HandleFunc("/users/{user_id}", controllers.User.GetUserByID).Methods(http.MethodGet)
	router.HandleFunc("/transactions/{user_id}", controllers.Transaction.GetAllTransactionsByUserID).Methods(http.MethodGet)
	router.HandleFunc("/process-record/{user_id}", controllers.Transaction.ProcessRecord).Methods(http.MethodPost)

	return router
}

func validateSourceHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !model.IsArrayContainsValue(enum.SourceTypes, r.Header.Get("Source-Type")) {
			controller.StatusBadRequest(w, "Source-Type header is not valid")
			return
		}

		next.ServeHTTP(w, r)
	})
}
