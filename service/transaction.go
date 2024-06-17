package service

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/entain-test-task/model"
	"github.com/entain-test-task/model/enum"
	requestmodel "github.com/entain-test-task/model/request"
	responsemodel "github.com/entain-test-task/model/response"
	"github.com/entain-test-task/repository"
	"github.com/go-openapi/strfmt"
	"github.com/pkg/errors"
)

type Transaction struct {
	transactionRepository *repository.Transaction
	userRepository        *repository.User
}

func NewTransaction(
	transactionRepository *repository.Transaction,
	userRepository *repository.User,
) *Transaction {
	return &Transaction{
		transactionRepository: transactionRepository,
		userRepository:        userRepository,
	}
}

func (service *Transaction) GetAllTransactionsByUserID(userID strfmt.UUID4) (*responsemodel.GetAllTransactionsByUserIDResponse, error) {
	transactions, err := service.transactionRepository.SelectTransactionsByUserID(userID)
	if err != nil {
		return nil, err
	}

	return &responsemodel.GetAllTransactionsByUserIDResponse{
		Transactions: transactions,
	}, nil
}

func (service *Transaction) ProcessRecord(userID strfmt.UUID4, processRecordRequest requestmodel.ProcessRecordRequest) (*responsemodel.ProcessRecordResponse, error) {
	var amount float64

	switch processRecordRequest.State {
	case enum.RecordStateWin:
		amount = processRecordRequest.Amount
	case enum.RecordStateLose:
		amount = -processRecordRequest.Amount
	}

	tx, err := service.transactionRepository.Begin()
	if err != nil {
		return nil, errors.Wrap(err, "beginning transaction failed")
	}

	defer service.rollbackTransaction(tx)

	// Check if the user has enough balance
	user, err := service.userRepository.SelectUser(userID)
	if err != nil {
		if err.Error() == model.ErrUserNotFound().Error() {
			return nil, model.ErrUserNotFound()
		}

		return nil, errors.Wrap(err, "selecting user failed")
	}

	if user.Balance+amount < 0 {
		return nil, model.ErrInsufficientBalance()
	}

	// Insert the transaction record
	transaction := model.Transaction{
		ID:     processRecordRequest.TransactionID,
		UserID: userID,
		Amount: amount,
	}

	if err := service.transactionRepository.Insert(tx, transaction); err != nil {
		if err.Error() == model.ErrTransactionAlreadyExists().Error() {
			return nil, model.ErrTransactionAlreadyExists()
		}

		return nil, errors.Wrap(err, "inserting transaction failed")
	}

	// Update the user balance
	if err := service.userRepository.UpdateUserBalance(tx, userID, amount); err != nil {
		return nil, errors.Wrap(err, "updating user balance failed")
	}

	if err := tx.Commit(); err != nil {
		return nil, errors.Wrap(err, "committing transaction failed")
	}

	return &responsemodel.ProcessRecordResponse{
		Message: http.StatusText(http.StatusOK),
	}, nil
}

func (service *Transaction) CancelLatestOddTransactionRecords(numberOfLatestRecords int) {
	transactions, err := service.transactionRepository.SelectLatestOddRecordTransactions(numberOfLatestRecords)
	if err != nil {
		log.Printf("selecting latest odd record transactions failed. %v", err)
		return
	}

	if len(transactions) == 0 {
		log.Printf("no records to cancel")
		return
	}

	for _, transaction := range transactions {
		if err = service.cancelTransactionRecord(transaction); err != nil {
			log.Printf("beginning transaction failed. %v", err)
			return
		}
	}

	log.Printf("successfully cancelled %d records", len(transactions))
}

func (service *Transaction) cancelTransactionRecord(transaction model.Transaction) error {
	tx, err := service.transactionRepository.Begin()
	if err != nil {
		return errors.Wrap(err, "beginning transaction failed")
	}

	defer service.rollbackTransaction(tx)

	// Cancel the transaction
	today := time.Now()
	transaction.CanceledAt = &today
	if err := service.transactionRepository.Update(tx, transaction); err != nil {
		return errors.Wrap(err, "canceling transaction failed")
	}

	// Refund the user balance
	if err := service.userRepository.UpdateUserBalance(tx, transaction.UserID, -transaction.Amount); err != nil {
		return errors.Wrap(err, "updating user balance failed")
	}

	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "committing transaction failed")
	}

	return nil
}

func (service *Transaction) rollbackTransaction(tx *sql.Tx) {
	if err := tx.Rollback(); err != nil {
		if err.Error() != "sql: transaction has already been committed or rolled back" {
			log.Printf("rolling back transaction failed. %v", err)
		}
	}
}
