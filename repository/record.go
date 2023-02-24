package repository

import (
	"log"
	"time"

	"github.com/entain-test-task/model"
	requestmodel "github.com/entain-test-task/model/request"
	"github.com/go-openapi/strfmt"
	"github.com/pkg/errors"
)

func ProcessRecord(userID strfmt.UUID4, processRecordRequest requestmodel.ProcessRecordRequest) error {
	var amount float64

	switch processRecordRequest.State {
	case model.RecordStateWin:
		amount = processRecordRequest.Amount
	case model.RecordStateLose:
		amount = -processRecordRequest.Amount
	}

	tx, err := DB.Begin()
	if err != nil {
		return errors.Wrap(err, "failed to begin transaction")
	}

	defer func() {
		if err := tx.Rollback(); err != nil {
			log.Printf("failed to rollback transaction: %v", err)
		}
	}()

	// Check if the user has enough balance
	var balance float64
	if err := tx.QueryRow(`
		SELECT
			balance
		FROM
			users
		WHERE
			id = $1
	`,
		userID,
	).Scan(&balance); err != nil {
		return errors.Wrap(err, "failed to get user balance")
	}

	if balance+amount < 0 {
		return ErrInsufficientBalance()
	}

	// Insert the transaction
	if _, err := tx.Exec(`
		INSERT INTO
			transaction (id, user_id, amount, created_at)
		VALUES
			($1, $2, $3, $4)
	`,
		processRecordRequest.TransactionID,
		userID,
		amount,
		time.Now(),
	); err != nil {
		if err.Error() == "pq: duplicate key value violates unique constraint \"transaction_pkey\"" {
			return ErrTransactionAlreadyExists()
		}

		return errors.Wrap(err, "failed to insert transaction")
	}

	// Update the user balance
	if _, err := tx.Exec(`
		UPDATE
			users
		SET
			balance = balance + $1,
			updated_at = $2
		WHERE
			id = $3
	`,
		amount,
		time.Now(),
		userID,
	); err != nil {
		return errors.Wrap(err, "failed to update user balance")
	}

	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "failed to commit transaction")
	}

	return nil
}

func ErrTransactionAlreadyExists() error {
	return errors.New("transaction already exists")
}

func ErrInsufficientBalance() error {
	return errors.New("insufficient balance")
}
