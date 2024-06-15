package controller

import (
	"github.com/entain-test-task/repository"

	_ "github.com/lib/pq"
)

type Controllers struct {
	User *User
}

func NewControllers(
	repository *repository.Store,
) *Controllers {
	return &Controllers{
		User: NewUser(repository),
	}
}
