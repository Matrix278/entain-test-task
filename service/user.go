package service

import (
	"github.com/entain-test-task/model"
	responsemodel "github.com/entain-test-task/model/response"
	"github.com/entain-test-task/repository"
	"github.com/go-openapi/strfmt"
)

type User struct {
	userRepository *repository.User
}

func NewUser(userRepository *repository.User) *User {
	return &User{
		userRepository: userRepository,
	}
}

func (service *User) GetAllUsers() (*responsemodel.GetAllUsersResponse, error) {
	users, err := service.userRepository.SelectUsers()
	if err != nil {
		return nil, err
	}

	return &responsemodel.GetAllUsersResponse{
		Users: users,
	}, nil
}

func (service *User) GetUser(userID strfmt.UUID4) (*responsemodel.GetUserResponse, error) {
	user, err := service.userRepository.SelectUser(userID)
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
