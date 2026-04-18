package main

import (
	"BlockMe/config"
	bm_cron "BlockMe/cron"
	"BlockMe/handlers"
	"BlockMe/logger"
	"BlockMe/middleware"
	"BlockMe/to"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/netresearch/go-cron"
)

func main() {
	fmt.Println("Setting up...")
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	err := logger.Init()
	if err != nil {
		log.Fatal(err)
	}

	err = config.InitConfig()
	if err != nil {
		fmt.Printf("Error initializing config: %s", err.Error())
		os.Exit(1)
	}

	err = to.Setup()
	if err != nil {
		fmt.Printf("Error setting up the database connection: %s\n", err.Error())
		os.Exit(1)
	}

	var c *cron.Cron
	if config.Env.IAmALittleBitch {
		c = bm_cron.InitCronJobs()
	}

	router := mux.NewRouter()
	router.Use(
		middleware.RealIpMiddleware,
		middleware.LogRequestMiddleware,
		middleware.BlockCheckMiddleware,
	)

	router.HandleFunc("/", handlers.Home).Methods("GET")
	router.HandleFunc("/blockme", handlers.BlockThem).Methods("POST")
	router.HandleFunc("/reset", handlers.ResetThem).Methods("GET")

	port := 8080
	server := &http.Server{
		Handler:      router,
		Addr:         fmt.Sprintf(":%d", port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	go func() {
		fmt.Printf("\nStarting server on port %d\n", port)
		if err = server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("The web server failed to start... %s\n", err.Error())
		}
	}()

	<-sigs
	fmt.Printf("\nReceived termination signal. Starting graceful shutdown.\n")

	if c != nil {
		fmt.Println("Stopping the cron jobs...")
		c.Stop()
	}

	fmt.Println("Stopping web server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		log.Fatalf("The web server failed to shut down gracefully... %s\n", err.Error())
	}

	fmt.Println("Closing the database connection...")
	err = to.DB.Close()
	if err != nil {
		fmt.Printf("Error closing the database connection: %s\n", err.Error())
	}

	fmt.Println("\nServer gracefully stopped")
}
