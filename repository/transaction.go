package repository

import (
	"github.com/entain-test-task/model"
	"github.com/go-openapi/strfmt"
	"github.com/pkg/errors"
)

func (repository *Store) GetAllTransactionsByUserID(userID strfmt.UUID4) ([]model.Transaction, error) {
	var transactions []model.Transaction

	rows, err := repository.db.Query(`
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
