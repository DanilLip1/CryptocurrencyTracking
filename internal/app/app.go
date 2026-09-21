package app

import (
	"context"
	"cryptocurrency/deploy/config"
	"cryptocurrency/internal/adapters/client/coingecko"
	"cryptocurrency/internal/adapters/repository/postgres"
	"cryptocurrency/internal/cases"
	"cryptocurrency/internal/ports/http/public"
	"log"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
)

func NewApp(ctx context.Context, cfg *config.Config) error {
	log.Println("initializing application")

	repository, err := postgres.NewRepository(ctx, cfg.PostgresURL)
	if err != nil {
		log.Printf("app: failed to create repository: %v", err)
		return errors.Wrap(err, "app: create postgres repository")
	}
	defer repository.Close()
	log.Println("app: repository created")

	migration, err := migrate.New("file:///app/deploy/migration/postgres", cfg.PostgresURL)
	if err != nil {
		log.Printf("app: failed to create migration : %v", err)
		return errors.Wrap(err, "app: create migration ")
	}
	defer migration.Close()
	if err := migration.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return errors.Wrap(err, "app: migration up")
	}
	log.Println("app: migration up")

	provider, err := coingecko.NewClient(cfg.CoinGeckoApiKey, cfg.CoinGeckoBaseURL)
	if err != nil {
		log.Printf("app: failed to create provider: %v", err)
		return errors.Wrap(err, "app: create provider coingecko")
	}
	log.Println("app: provider created")

	service, err := cases.NewCurrencyService(repository, provider)
	if err != nil {
		log.Printf("app: filed to create service: %v", err)
		return errors.Wrap(err, "app: create service")
	}
	log.Println("app: service created")

	server, err := public.NewServer(service, cfg.HTTPAddress)
	if err != nil {
		log.Printf("app: failed to create server: %v", err)
		return errors.Wrap(err, "app: create server")
	}
	log.Println("app: server created")

	scheduler := cron.New()
	_, err = scheduler.AddFunc("@every "+cfg.CronUpdateInterval, func() {
		log.Println("app: cron update prices started")
		if err := service.UpdatePrices(context.Background()); err != nil {
			log.Printf("app: cron update prices failed: %v", err)
			return
		}
		log.Println("app: cron update prices completed")
	})
	if err != nil {
		log.Printf("app: failed to add cron update job: %v", err)
		return errors.Wrap(err, "app: add cron job")
	}

	scheduler.Start()
	log.Println("app: scheduler started")

	serverErrors := make(chan error, 1)
	go func() {
		log.Println("app: starting http server")

		if err := server.Run(); err != nil {
			serverErrors <- err
		}
	}()

	select {
	case err := <-serverErrors:
		log.Println("app: server cron job finished with error:", err)
		scheduler.Stop()
		return errors.Wrap(err, "app: http server failed")
	case <-ctx.Done():
		log.Println("app: graceful shutdown started")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("app: http server shutdown failed: %v", err)
		return errors.Wrap(err, "app: shutdown HTTP server")
	}
	log.Println("app: http server stopped")

	schedulerCtx := scheduler.Stop()
	<-schedulerCtx.Done()
	log.Println("app: scheduler stopped")

	log.Println("app: graceful shutdown completed")
	return nil
}
