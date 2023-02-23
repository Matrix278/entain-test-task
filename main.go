package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/entain-test-task/middleware"
	"github.com/entain-test-task/repository"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println(fmt.Errorf("error loading .env file: %v", err))
	}

	db := repository.InitDB()
	defer db.Close()

	router := middleware.Router()

	fmt.Println("Starting server on the port 8080...")

	log.Fatal(http.ListenAndServe(":8080", router))
}
