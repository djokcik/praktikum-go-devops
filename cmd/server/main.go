package main

import (
	"context"
	"flag"
	"github.com/caarlos0/env/v6"
	"github.com/djokcik/praktikum-go-devops/internal/server"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}

	var cfg server.Config

	parseEnv(&cfg)
	parseFlags(&cfg)

	log.Println(cfg)

	mux := chi.NewMux()

	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)

	makeMetricRoutes(ctx, wg, mux, &cfg)

	go func() {
		log.Fatal(http.ListenAndServe(cfg.Address, mux))
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	<-quit
	cancel()
	log.Println("Shutdown Server ...")
	wg.Wait()
}

func parseEnv(cfg *server.Config) {
	err := env.Parse(cfg)
	if err != nil {
		log.Fatal(err)
	}
}

func parseFlags(cfg *server.Config) {
	flag.Func("a", "Server address", func(address string) error {
		if _, ok := os.LookupEnv("ADDRESS"); ok {
			return nil
		}

		cfg.Address = address
		return nil
	})

	flag.Func("r", "Restore", func(restorePlan string) error {
		if _, ok := os.LookupEnv("RESTORE"); ok {
			return nil
		}

		restore, err := strconv.ParseBool(restorePlan)
		if err != nil {
			return err
		}

		cfg.Restore = restore
		return nil
	})

	flag.Func("i", "Store save interval", func(storeIntervalPlan string) error {
		if _, ok := os.LookupEnv("STORE_INTERVAL"); ok {
			return nil
		}

		storeInterval, err := time.ParseDuration(storeIntervalPlan)
		if err != nil {
			return err
		}

		cfg.StoreInterval = storeInterval
		return nil
	})

	flag.Func("f", "Store file", func(storeFile string) error {
		if _, ok := os.LookupEnv("STORE_FILE"); ok {
			return nil
		}

		cfg.StoreFile = storeFile
		return nil
	})

	flag.Parse()
}
