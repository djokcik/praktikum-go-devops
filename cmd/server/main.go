package main

import (
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

func main() {
	mux := chi.NewMux()

	makeMetricRoutes(mux)

	log.Fatal(http.ListenAndServe(":8080", mux))
}
