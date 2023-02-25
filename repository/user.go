package repository

import (
	"github.com/entain-test-task/model"
	"github.com/go-openapi/strfmt"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

func GetAllUsers() ([]model.User, error) {
	var users []model.User

	rows, err := DB.Query(`
		SELECT
			*
		FROM
			users
	`)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get all users")
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

func GetUser(userID strfmt.UUID4) (*model.User, error) {
	var user model.User

	err := DB.QueryRow(`
		SELECT
			*
		FROM
			users
		WHERE
			id = $1
	`,
		userID,
	).Scan(&user.ID, &user.Balance, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, ErrUserNotFound()
		}

		return nil, errors.Wrap(err, "failed to get user")
	}

	return &user, nil
}

func ErrUserNotFound() error {
	return errors.New("user not found")
}
