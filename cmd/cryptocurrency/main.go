package main

import (
	"context"
	"cryptocurrency/deploy/config"
	_ "cryptocurrency/docs"
	"cryptocurrency/internal/app"
	"log"
	"os/signal"
	"syscall"
)

// @title Cryptocurrency Tracking API
// @version 1.0
// @description API for tracking cryptocurrency prices.
// @host localhost:8080
// @BasePath /api/v1
func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	log.Println("main: application started")

	cfg, err := config.LoadConfig("deploy/config/config.yaml")
	if err != nil {
		log.Printf("main: failed to load config: %v", err)
		return
	}

	if err := app.NewApp(ctx, cfg); err != nil {
		log.Printf("main: failed to start app: %v", err)
		return
	}

	log.Println("main: application stopped")
}
