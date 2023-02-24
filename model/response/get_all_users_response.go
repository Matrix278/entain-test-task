package model

import "github.com/entain-test-task/model"

type GetAllUsersResponse struct {
	Users []model.User `json:"users"`
}
