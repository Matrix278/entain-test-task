package repository

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/entain-test-task/configuration"
	_ "github.com/lib/pq"
)

type Store struct {
	db *sql.DB
}

func NewStore(cfg configuration.Config) *Store {
	db, err := sql.Open("postgres", cfg.PostgresDatabaseURL)
	if err != nil {
		log.Fatalf("unable to connect to database. %v", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("unable to connect to database. %v", err)
	}

	fmt.Println("Database successfully connected!")

	return &Store{
		db: db,
	}
}

func (repository *Store) Close() {
	repository.db.Close()
}
