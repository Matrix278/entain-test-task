package model

import "github.com/entain-test-task/model"

type GetAllTransactionsByUserIDResponse struct {
	Transactions []model.Transaction `json:"transactions"`
}
