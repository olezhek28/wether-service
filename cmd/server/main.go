package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-co-op/gocron/v2"
	"github.com/olezhek28/wether-service/internal/client/http/geocoding"
	"github.com/olezhek28/wether-service/internal/client/http/open_meteo"
)

const (
	httpPort = ":3000"
	city     = "moscow"
)

type Reading struct {
	Timestamp   time.Time
	Temperature float64
}

type Storage struct {
	data map[string][]Reading
	mu   sync.RWMutex
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	storage := &Storage{
		data: make(map[string][]Reading),
	}

	r.Get("/{city}", func(w http.ResponseWriter, r *http.Request) {
		cityName := chi.URLParam(r, "city")

		fmt.Printf("Requested city: %s\n", cityName)

		storage.mu.RLock()
		defer storage.mu.RUnlock()

		reading, ok := storage.data[cityName]
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("not found"))
			return
		}

		raw, err := json.Marshal(reading)
		if err != nil {
			log.Println(err)
		}

		_, err = w.Write(raw)
		if err != nil {
			log.Println(err)
		}
	})

	s, err := gocron.NewScheduler()
	if err != nil {
		panic(err)
	}

	jobs, err := initJobs(s, storage)
	if err != nil {
		panic(err)
	}

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()

		fmt.Println("starting server on port " + httpPort)
		err = http.ListenAndServe(httpPort, r)
		if err != nil {
			panic(err)
		}
	}()

	go func() {
		defer wg.Done()

		fmt.Printf("starting job: %v\n", jobs[0].ID())
		s.Start()
	}()

	wg.Wait()
}

func initJobs(scheduler gocron.Scheduler, storage *Storage) ([]gocron.Job, error) {
	httpClient := &http.Client{
		Timeout: time.Second * 10,
	}
	geocodingClient := geocoding.NewClient(httpClient)
	openMeteoClient := open_meteo.NewClient(httpClient)

	j, err := scheduler.NewJob(
		gocron.DurationJob(
			10*time.Second,
		),
		gocron.NewTask(
			func() {
				geoRes, err := geocodingClient.GetCoords(city)
				if err != nil {
					log.Println(err)
					return
				}

				openMetRes, err := openMeteoClient.GetTemperature(geoRes.Latitude, geoRes.Longitude)
				if err != nil {
					log.Println(err)
					return
				}

				storage.mu.Lock()
				defer storage.mu.Unlock()

				timestamp, err := time.Parse("2006-01-02T15:04", openMetRes.Current.Time)
				if err != nil {
					log.Println(err)
					return
				}

				storage.data[city] = append(storage.data[city], Reading{
					Timestamp:   timestamp,
					Temperature: openMetRes.Current.Temperature2m,
				})

				fmt.Printf("%v updated data for city: %s\n", time.Now(), city)
			},
		),
	)
	if err != nil {
		return nil, err
	}

	return []gocron.Job{j}, nil
}
