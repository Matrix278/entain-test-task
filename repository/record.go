package repository

import (
	"time"

	"github.com/entain-test-task/model"
	requestmodel "github.com/entain-test-task/model/request"
	"github.com/go-openapi/strfmt"
	"github.com/pkg/errors"
)

func ProcessRecord(userID strfmt.UUID4, processRecordRequest requestmodel.ProcessRecordRequest) error {
	var amountProcessIdentifier string

	switch processRecordRequest.State {
	case model.RecordStateWin:
		amountProcessIdentifier = "+"
	case model.RecordStateLose:
		amountProcessIdentifier = "-"
	}

	tx, err := DB.Begin()
	if err != nil {
		return errors.Wrap(err, "failed to begin transaction")
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
		processRecordRequest.Amount,
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
			balance = balance `+amountProcessIdentifier+` $1,
			updated_at = $2
		WHERE
			id = $3
	`,
		processRecordRequest.Amount,
		time.Now(),
		userID,
	); err != nil {
		if err := tx.Rollback(); err != nil {
			return errors.Wrap(err, "failed to rollback transaction")
		}

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
