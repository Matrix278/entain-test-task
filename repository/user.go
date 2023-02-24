package repository

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/entain-test-task/model"
	"github.com/go-openapi/strfmt"
	_ "github.com/lib/pq"
)

var DB *sql.DB

// InitDB will initialize the database
func InitDB() *sql.DB {
	psqlInfo := fmt.Sprintf(
		"host=%s port=%s user=%s "+"password=%s dbname=%s sslmode=disable",
		os.Getenv("POSTGRES_DB_HOST"),
		os.Getenv("POSTGRES_DB_PORT"),
		os.Getenv("POSTGRES_DB_USER"),
		os.Getenv("POSTGRES_DB_PASS"),
		os.Getenv("POSTGRES_DB_NAME"),
	)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("Unable to connect to database. %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Unable to connect to database. %v", err)
	}

	fmt.Println("DB successfully connected!")
	DB = db

	return db
}

// GetAllUsers will return all the users in the database
func GetAllUsers() ([]model.User, error) {
	var users []model.User

	rows, err := DB.Query(`
		SELECT
			*
		FROM
			users
	`)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var user model.User
		err := rows.Scan(&user.ID, &user.Balance, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

// GetUser will return a user by id
func GetUser(userID strfmt.UUID4) (model.User, error) {
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
		return model.User{}, err
	}

	return user, nil
}
