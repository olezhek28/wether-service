package app

import (
	"context"

	"github.com/go-chi/chi/v5"
	"github.com/go-co-op/gocron/v2"
	"github.com/olezhek28/wether-service/internal/config"
	"github.com/olezhek28/wether-service/internal/cron"
	"github.com/olezhek28/wether-service/internal/handlers"
	"github.com/olezhek28/wether-service/internal/http"
	"github.com/olezhek28/wether-service/internal/services"
	"github.com/olezhek28/wether-service/internal/storage"
	"github.com/olezhek28/wether-service/internal/storage/postgres"
)

type App struct {
	Server *http.Server
	Cron   gocron.Scheduler
}

func New(config *config.Config) *App {
	ctx := context.Background()

	r := chi.NewRouter()

	scheduler, err := gocron.NewScheduler()
	if err != nil {
		panic(err)
	}

	srv := http.NewServer(ctx, config.Port, config.Host, r)

	postgres := postgres.New(ctx, config)

	weatherDB := storage.New(postgres)

	service := services.New(weatherDB, weatherDB)

	h := handlers.New(r, service)
	h.Init()

	c := cron.New(scheduler, service)
	c.Init(ctx)

	return &App{
		Server: srv,
		Cron:   scheduler,
	}
}
