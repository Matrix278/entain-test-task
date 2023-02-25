package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

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

	nMinutes, err := strconv.Atoi(os.Getenv("CANCEL_ODD_RECORDS_MINUTES_INTERVAL"))
	if err != nil {
		log.Fatal("CANCEL_ODD_RECORDS_MINUTES_INTERVAL must be an integer")
	}

	go func() {
		for {
			time.Sleep(time.Duration(nMinutes) * time.Minute)
			middleware.CancelLatestOddTransactionRecords(10)
		}
	}()

	fmt.Println("Starting server on the port " + os.Getenv("SERVER_PORT") + "...")

	log.Fatal(http.ListenAndServe(":"+os.Getenv("SERVER_PORT"), router))
}
