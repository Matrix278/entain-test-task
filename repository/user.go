package repository

import (
	"database/sql"
	"time"

	"github.com/entain-test-task/model"
	"github.com/go-openapi/strfmt"
	"github.com/pkg/errors"
)

type User struct {
	store *Store
}

func NewUser(store *Store) *User {
	return &User{
		store: store,
	}
}

func (repository *User) Begin() (*sql.Tx, error) {
	return repository.store.db.Begin()
}

func (repository *User) SelectUsers() ([]model.User, error) {
	var users []model.User

	rows, err := repository.store.db.Query(`
		SELECT
			*
		FROM
			users
	`)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select users")
	}

	for rows.Next() {
		var user model.User
		err := rows.Scan(
			&user.ID,
			&user.Balance,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan user")
		}

		users = append(users, user)
	}

	return users, nil
}

func (repository *User) SelectUser(userID strfmt.UUID4) (*model.User, error) {
	var user model.User

	if err := repository.store.db.QueryRow(`
		SELECT
			*
		FROM
			users
		WHERE
			id = $1
	`,
		userID,
	).Scan(&user.ID, &user.Balance, &user.CreatedAt, &user.UpdatedAt); err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, model.ErrUserNotFound()
		}

		return nil, errors.Wrap(err, "failed to select user")
	}

	return &user, nil
}

func (repository *User) UpdateUserBalance(tx *sql.Tx, userID strfmt.UUID4, amount float64) error {
	if tx == nil {
		return errors.New("tx is nil")
	}

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

	return nil
}
