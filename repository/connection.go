package repository

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

type Store struct {
	db *sql.DB
}

func NewStore() *Store {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("POSTGRES_DB_HOST"),
		os.Getenv("POSTGRES_DB_PORT"),
		os.Getenv("POSTGRES_DB_USER"),
		os.Getenv("POSTGRES_DB_PASS"),
		os.Getenv("POSTGRES_DB_NAME"),
	)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("unable to connect to database. %v", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("unable to connect to database. %v", err)
	}

	fmt.Println("DB successfully connected!")

	return &Store{
		db: db,
	}
}

func (repository *Store) Close() {
	repository.db.Close()
}
