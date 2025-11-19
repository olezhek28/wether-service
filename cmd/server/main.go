package main

import (
	"sync"

	"github.com/olezhek28/wether-service/internal/app"
	"github.com/olezhek28/wether-service/internal/config"
)

func main() {
	var wg sync.WaitGroup
	cfg := config.MustLoad()

	app := app.New(cfg)

	wg.Add(2)
	go func() {
		defer wg.Done()
		app.Server.MustRun()
	}()

	go func() {
		defer wg.Done()
		app.Cron.Start()
	}()

	wg.Wait()
}
