package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const httpPort = ":3000"

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("привет!"))
		if err != nil {
			log.Println(err)
		}

	})

	err := http.ListenAndServe(httpPort, r)
	if err != nil {
		panic(err)
	}
}
