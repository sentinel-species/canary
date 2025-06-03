package main

import (
	"canary/internal/config"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"canary/cmd/server"
)

func main() {
	err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	errChan := make(chan error, 1)
	go func() {
		errChan <- server.Start()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errChan:
		log.Fatalf("Server error: %v", err)
	case sig := <-quit:
		log.Printf("Received signal: %v", sig)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Stop(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited properly")
}
