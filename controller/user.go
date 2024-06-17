package controller

import (
	"log"
	"net/http"

	"github.com/entain-test-task/model"
	"github.com/entain-test-task/service"
	"github.com/go-openapi/strfmt"
	"github.com/gorilla/mux"
)

type User struct {
	service service.IUser
}

func NewUser(service service.IUser) *User {
	return &User{
		service: service,
	}
}

func (controller *User) GetAllUsers(w http.ResponseWriter, _ *http.Request) {
	getAllUsersResponse, err := controller.service.GetAllUsers()
	if err != nil {
		log.Printf("getting all users failed. %v", err)
		StatusInternalServerError(w)
		return
	}

	StatusOK(w, getAllUsersResponse)
}

func (controller *User) GetUserByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID := params["user_id"]

	if !strfmt.IsUUID4(userID) {
		StatusBadRequest(w, "user_id is not a valid UUID4")
		return
	}

	user, err := controller.service.GetUser(strfmt.UUID4(userID))
	if err != nil {
		if err.Error() == model.ErrUserNotFound().Error() {
			StatusUnprocessableEntity(w, model.ErrUserNotFound().Error())
			return
		}

		log.Printf("getting user by ID failed. %v", err)
		StatusInternalServerError(w)
		return
	}

	StatusOK(w, user)
}
