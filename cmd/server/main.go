package main

import (
	"fmt"
	"github.com/Jokcik/praktikum-go-devops/internal/server/handlers"
	"github.com/Jokcik/praktikum-go-devops/internal/server/storage"
	"github.com/Jokcik/praktikum-go-devops/internal/server/storage/model"
	"log"
	"net/http"
)

func main() {
	serveHandler := http.NewServeMux()

	rr := storage.NewRepositoryRegistry(new(model.Database), &storage.MetricRepository{})
	metricRepository, err := rr.Repository("MetricRepository")
	if err != nil {
		fmt.Println("Error provide repository 'MetricRepository'")
	}

	serveHandler.HandleFunc("/update/gauge/", handlers.GaugeHandler(metricRepository))
	serveHandler.HandleFunc("/update/counter/", handlers.CounterHandler(metricRepository))
	serveHandler.HandleFunc("/update/", handlers.NotImplementedHandler)

	log.Fatal(http.ListenAndServe(":8080", serveHandler))
}
