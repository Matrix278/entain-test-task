package controller

import (
	"github.com/entain-test-task/repository"
	"github.com/entain-test-task/service"

	_ "github.com/lib/pq"
)

type Controllers struct {
	User        *User
	Transaction *Transaction
}

func NewControllers(
	repository *repository.Store,
) *Controllers {
	return &Controllers{
		User:        NewUser(service.NewUser(repository)),
		Transaction: NewTransaction(service.NewTransaction(repository)),
	}
}
