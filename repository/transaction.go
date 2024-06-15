package repository

import (
	"database/sql"
	"time"

	"github.com/entain-test-task/model"
	"github.com/go-openapi/strfmt"
	"github.com/pkg/errors"
)

type Transaction struct {
	store *Store
}

func NewTransaction(store *Store) *Transaction {
	return &Transaction{
		store: store,
	}
}

func (repository *Transaction) Begin() (*sql.Tx, error) {
	return repository.store.db.Begin()
}

func (repository *Transaction) SelectTransactionsByUserID(userID strfmt.UUID4) ([]model.Transaction, error) {
	var transactions []model.Transaction

	rows, err := repository.store.db.Query(`
		SELECT
			*
		FROM
			transaction
		WHERE
			user_id = $1
	`,
		userID,
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get all transactions")
	}

	for rows.Next() {
		var transaction model.Transaction
		err := rows.Scan(
			&transaction.ID,
			&transaction.UserID,
			&transaction.Amount,
			&transaction.CreatedAt,
			&transaction.CanceledAt,
		)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan transaction")
		}

		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

func (repository *Transaction) Insert(tx *sql.Tx, transaction model.Transaction) error {
	if tx == nil {
		return errors.New("tx is nil")
	}

	if _, err := tx.Exec(`
		INSERT INTO
			transaction (id, user_id, amount, created_at)
		VALUES
			($1, $2, $3, $4)
	`,
		transaction.ID,
		transaction.UserID,
		transaction.Amount,
		time.Now(),
	); err != nil {
		if err.Error() == "pq: duplicate key value violates unique constraint \"transaction_pkey\"" {
			return model.ErrTransactionAlreadyExists()
		}

		return errors.Wrap(err, "failed to insert transaction")
	}

	return nil
}

func (repository *Transaction) Update(tx *sql.Tx, transaction model.Transaction) error {
	if tx == nil {
		return errors.New("tx is nil")
	}

	if _, err := tx.Exec(`
		UPDATE
			transaction
		SET
			user_id = $1,
			amount = $2,
			canceled_at = $3
		WHERE
			id = $4
	`,
		transaction.UserID,
		transaction.Amount,
		transaction.CanceledAt,
		transaction.ID,
	); err != nil {
		return errors.Wrap(err, "failed to update transaction")
	}

	return nil
}

func (repository *Transaction) SelectLatestOddRecordTransactions(numberOfTransactionRecords int) ([]model.Transaction, error) {
	var transactions []model.Transaction

	rows, err := repository.store.db.Query(`
		SELECT
			*
		FROM
			transaction
		WHERE
			(
				mod(amount, 2) = 1
			OR
				mod(amount, 2) = -1
			)
		AND
			canceled_at IS NULL
		ORDER BY
			created_at DESC
		LIMIT
			$1
	`,
		numberOfTransactionRecords,
	)
	if err != nil {
		return transactions, errors.Wrap(err, "failed to get latest odd records")
	}

	defer rows.Close()

	for rows.Next() {
		var transaction model.Transaction

		if err := rows.Scan(
			&transaction.ID,
			&transaction.UserID,
			&transaction.Amount,
			&transaction.CreatedAt,
			&transaction.CanceledAt,
		); err != nil {
			return transactions, errors.Wrap(err, "failed to scan transaction")
		}

		transactions = append(transactions, transaction)
	}

	return transactions, nil
}
