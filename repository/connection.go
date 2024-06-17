package repository

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/entain-test-task/configuration"
	_ "github.com/lib/pq" // nolint: revive
)

type Store struct {
	db *sql.DB
}

func NewStore(cfg *configuration.Config) *Store {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.PostgresHost,
		cfg.PostgresPort,
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresName,
	)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("connecting to database failed. %v", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("connecting to database failed. %v", err)
	}

	log.Println("Database successfully connected!")

	return &Store{
		db: db,
	}
}

func (repository *Store) Close() {
	if err := repository.db.Close(); err != nil {
		log.Fatalf("closing database connection failed. %v", err)
	}
}
