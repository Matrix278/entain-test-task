package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/entain-test-task/configuration"
	"github.com/entain-test-task/controller"
	"github.com/entain-test-task/middleware"
	"github.com/entain-test-task/repository"
)

func main() {
	cfg, err := configuration.Load()
	if err != nil {
		log.Fatal(fmt.Errorf("error loading configuration: %v", err))
	}

	store := repository.NewStore(cfg)
	defer store.Close()

	controllers := controller.NewControllers(cfg, store)

	router := middleware.Router(controllers)

	go func() {
		for {
			time.Sleep(time.Duration(cfg.CancelOddRecordsMinutesInterval) * time.Minute)
			controllers.Transaction.CancelLatestOddTransactionRecords()
		}
	}()

	fmt.Println("Starting server on the port " + cfg.ServerPort + "...")

	log.Fatal(http.ListenAndServe(":"+cfg.ServerPort, router))
}
