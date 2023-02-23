package middleware

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/entain-test-task/repository"

	_ "github.com/lib/pq"
)

func ProcessRecord(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// TODO: implement the logic to process the record
	users, err := repository.GetAllUsers()
	if err != nil {
		log.Fatalf("Unable to get all users. %v", err)
	}

	json.NewEncoder(w).Encode(users)
}
