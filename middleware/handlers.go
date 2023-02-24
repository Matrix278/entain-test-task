package middleware

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/entain-test-task/model"
	"github.com/entain-test-task/repository"
	"github.com/go-openapi/strfmt"
	"github.com/gorilla/mux"

	_ "github.com/lib/pq"
)

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := repository.GetAllUsers()
	if err != nil {
		log.Printf("Unable to get all users. %v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	getAllUsersResponse := model.GetAllUsersResponse{
		Users: users,
	}

	json.NewEncoder(w).Encode(getAllUsersResponse)
}

func GetUserByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID := params["id"]

	if !strfmt.IsUUID4(userID) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := repository.GetUser(strfmt.UUID4(userID))
	if err != nil {
		log.Printf("Unable to get user. %v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(user)
}
