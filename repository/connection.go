package repository

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

// InitDB will initialize the database
func InitDB() *sql.DB {
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

	err = db.Ping()
	if err != nil {
		log.Fatalf("unable to connect to database. %v", err)
	}

	fmt.Println("DB successfully connected!")
	DB = db

	return db
}
