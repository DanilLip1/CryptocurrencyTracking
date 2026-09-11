package main

import (
	"context"
	"cryptocurrency/deploy/config"
	_ "cryptocurrency/docs"
	"cryptocurrency/internal/app"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pkg/errors"
)

// @title Cryptocurrency Tracking API
// @version 1.0
// @description API for tracking cryptocurrency prices.
// @host localhost:8080
// @BasePath /api/v1
func main() {
	ctx := context.Background()

	cfg, err := config.LoadConfig("deploy/config/config.yaml")
	if err != nil {
		panic(err)
	}

	application, err := app.NewApp(ctx, cfg)
	if err != nil {
		panic(err)
	}

	go func() {
		if err := application.Start(); err != nil {
			panic(err)
		}
	}()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := application.Shutdown(shutdownCtx); err != nil {
		panic(errors.Wrap(err, "main: shutdown"))

	}
}
