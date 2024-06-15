package controller

import (
	"log"
	"net/http"

	responsemodel "github.com/entain-test-task/model/response"
	"github.com/entain-test-task/repository"
	"github.com/go-openapi/strfmt"
	"github.com/gorilla/mux"

	_ "github.com/lib/pq"
)

type User struct {
	repository *repository.Store
}

func NewUser(repository *repository.Store) *User {
	return &User{
		repository: repository,
	}
}

func (controller *User) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := controller.repository.GetAllUsers()
	if err != nil {
		log.Printf("unable to get all users. %v", err)
		StatusInternalServerError(w)
		return
	}

	StatusOK(w, responsemodel.GetAllUsersResponse{
		Users: users,
	})
}

func (controller *User) GetUserByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID := params["user_id"]

	if !strfmt.IsUUID4(userID) {
		StatusBadRequest(w, "user_id is not a valid UUID4")
		return
	}

	user, err := controller.repository.GetUser(strfmt.UUID4(userID))
	if err != nil {
		log.Printf("unable to get user. %v", err)

		if err.Error() == repository.ErrUserNotFound().Error() {
			StatusUnprocessableEntity(w, repository.ErrUserNotFound().Error())
			return
		}

		StatusInternalServerError(w)
		return
	}

	StatusOK(w, user)
}
