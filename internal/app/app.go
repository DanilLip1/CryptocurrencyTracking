package app

import (
	"context"
	"cryptocurrency/deploy/config"
	"cryptocurrency/internal/adapters/client/coingecko"
	"cryptocurrency/internal/adapters/repository/postgres"
	"cryptocurrency/internal/cases"
	"cryptocurrency/internal/ports/http/public"
	"log"

	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
)

type App struct {
	server     *public.Server
	repository *postgres.Repository
	cron       *cron.Cron
}

func NewApp(ctx context.Context, cfg *config.Config) (*App, error) {
	repository, err := postgres.NewRepository(ctx, cfg.Postgres.URL)
	if err != nil {
		return nil, errors.Wrap(err, "app: create postgres repository")
	}
	provider, err := coingecko.NewClient(cfg.CoinGecko.APIKey, cfg.CoinGecko.BaseURL)
	if err != nil {
		repository.Close()
		return nil, errors.Wrap(err, "app: create provider")
	}
	service, err := cases.NewCurrencyService(repository, provider)
	if err != nil {
		repository.Close()
		return nil, errors.Wrap(err, "app: create service")
	}
	server, err := public.NewServer(service, cfg.HTTP.Address)
	if err != nil {
		repository.Close()
		return nil, errors.Wrap(err, "app: create server")
	}

	scheduler := cron.New()
	_, err = scheduler.AddFunc("@every "+cfg.Cron.UpdateInterval, func() {
		if err := service.UpdatePrices(context.Background()); err != nil {
			log.Printf("cron update prices: %v", err)
			return
		}
	})
	if err != nil {
		repository.Close()
		return nil, errors.Wrap(err, "app: add cron job")
	}
	return &App{
		server:     server,
		repository: repository,
		cron:       scheduler,
	}, nil
}

func (app *App) Start() error {
	app.cron.Start()
	return app.server.Run()
}
func (app *App) Shutdown(ctx context.Context) error {
	app.cron.Stop()
	if err := app.server.Shutdown(ctx); err != nil {
		return errors.Wrap(err, "app: shutdown server")
	}
	app.repository.Close()
	return nil
}
