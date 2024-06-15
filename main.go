package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/entain-test-task/configuration"
	"github.com/entain-test-task/middleware"
	"github.com/entain-test-task/repository"
	"github.com/entain-test-task/service"
)

func main() {
	cfg, err := configuration.Load()
	if err != nil {
		log.Fatal(fmt.Errorf("error loading configuration: %v", err))
	}

	repo := repository.NewStore(cfg)
	defer repo.Close()

	router := middleware.Router(repo)

	transactionService := service.NewTransaction(repo)
	go func() {
		for {
			time.Sleep(time.Duration(cfg.CancelOddRecordsMinutesInterval) * time.Minute)
			transactionService.CancelLatestOddTransactionRecords(10)
		}
	}()

	fmt.Println("Starting server on the port " + cfg.ServerPort + "...")

	log.Fatal(http.ListenAndServe(":"+cfg.ServerPort, router))
}
