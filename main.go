package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/entain-test-task/configuration"
	"github.com/entain-test-task/controller"
	"github.com/entain-test-task/middleware"
	"github.com/entain-test-task/repository"
)

func main() {
	// Load the configuration
	cfg, err := configuration.Load()
	if err != nil {
		log.Fatal(fmt.Errorf("error loading configuration: %v", err))
	}

	// Initialize the db store
	store := repository.NewStore(cfg)
	defer store.Close()

	// Initialize the controllers
	controllers := controller.NewControllers(cfg, store)

	// Initialize the router
	router := middleware.Router(controllers)

	// Start a goroutine to cancel the latest odd transaction records.
	ctx, cancelCtx := context.WithCancel(context.Background())
	defer cancelCtx()

	var waitGroup sync.WaitGroup
	waitGroup.Add(1)
	go startCancelLatestOddTransactionRecords(ctx, cfg, controllers, &waitGroup)

	// Handle SIGINT and SIGTERM.
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	// Start the server in a separate goroutine.
	server := &http.Server{
		Addr:              ":" + cfg.ServerPort,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()

		log.Println("Starting server on the port " + cfg.ServerPort + "...")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	// Wait for a shutdown signal.
	<-ch

	// Cancel the context.
	cancelCtx()

	// Shutdown the server.
	if err := server.Shutdown(ctx); err != nil {
		log.Print(err)
	}

	// Wait for all goroutines to finish.
	waitGroup.Wait()

	log.Println("Shutting down...")
}

func startCancelLatestOddTransactionRecords(
	ctx context.Context,
	cfg *configuration.Config,
	controllers *controller.Controllers,
	waitGroup *sync.WaitGroup,
) {
	defer waitGroup.Done()

	if cfg.CancelOddRecordsMinutesInterval == 0 {
		log.Println("CancelOddRecordsMinutesInterval is set to 0. Canceling odd records is disabled.")
		return
	}

	ticker := time.NewTicker(time.Duration(cfg.CancelOddRecordsMinutesInterval) * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			controllers.Transaction.CancelLatestOddTransactionRecords()
		case <-ctx.Done():
			return
		}
	}
}
