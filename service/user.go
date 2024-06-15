package service

import (
	"github.com/entain-test-task/model"
	responsemodel "github.com/entain-test-task/model/response"
	"github.com/entain-test-task/repository"
	"github.com/go-openapi/strfmt"

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

func (service *User) GetAllUsers() (*responsemodel.GetAllUsersResponse, error) {
	users, err := service.repository.GetAllUsers()
	if err != nil {
		return nil, err
	}

	return &responsemodel.GetAllUsersResponse{
		Users: users,
	}, nil
}

func (service *User) GetUser(userID strfmt.UUID4) (*responsemodel.GetUserResponse, error) {
	user, err := service.repository.GetUser(userID)
	if err != nil {
		if err.Error() == model.ErrUserNotFound().Error() {
			return nil, model.ErrUserNotFound()
		}

		return nil, err
	}

	return &responsemodel.GetUserResponse{
		User: *user,
	}, nil
}
