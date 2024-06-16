package controller

import (
	"github.com/entain-test-task/configuration"
	"github.com/entain-test-task/repository"
	"github.com/entain-test-task/service"
)

type Controllers struct {
	User        *User
	Transaction *Transaction
}

func NewControllers(
	cfg *configuration.Config,
	store *repository.Store,
) *Controllers {
	userRepo := repository.NewUser(store)
	transactionRepo := repository.NewTransaction(store)

	userService := service.NewUser(userRepo)
	transactionService := service.NewTransaction(transactionRepo, userRepo)

	return &Controllers{
		User:        NewUser(userService),
		Transaction: NewTransaction(cfg, transactionService),
	}
}
