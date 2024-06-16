package repository

import (
	"database/sql"
	"log"

	"github.com/entain-test-task/configuration"
	_ "github.com/lib/pq" // nolint: revive
)

type Store struct {
	db *sql.DB
}

func NewStore(cfg *configuration.Config) *Store {
	db, err := sql.Open("postgres", cfg.PostgresDatabaseURL)
	if err != nil {
		log.Fatalf("unable to connect to database. %v", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("unable to connect to database. %v", err)
	}

	log.Println("Database successfully connected!")

	return &Store{
		db: db,
	}
}

func (repository *Store) Close() {
	if err := repository.db.Close(); err != nil {
		log.Fatalf("unable to close the database connection. %v", err)
	}
}
