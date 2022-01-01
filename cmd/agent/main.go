package main

import (
	"context"
	"flag"
	"github.com/caarlos0/env/v6"
	"github.com/djokcik/praktikum-go-devops/internal/agent"
	"github.com/djokcik/praktikum-go-devops/internal/helpers"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var cfg agent.Config

	parseEnv(&cfg)
	parseFlags(&cfg)

	metricAgent := agent.NewAgent(&cfg)

	go helpers.SetTicker(metricAgent.CollectMetrics(ctx), cfg.PollInterval)
	go helpers.SetTicker(metricAgent.SendToServer(ctx), cfg.ReportInterval)

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	<-quit
	log.Println("Shutdown Agent ...")
}

func parseEnv(cfg *agent.Config) {
	err := env.Parse(cfg)
	if err != nil {
		log.Fatal(err)
	}
}

func parseFlags(cfg *agent.Config) {
	flag.Func("a", "Server address", func(address string) error {
		if _, ok := os.LookupEnv("ADDRESS"); ok {
			return nil
		}

		cfg.Address = address
		return nil
	})

	flag.Func("r", "Report Interval", func(reportIntervalPlain string) error {
		if _, ok := os.LookupEnv("REPORT_INTERVAL"); ok {
			return nil
		}

		reportInterval, err := time.ParseDuration(reportIntervalPlain)
		if err != nil {
			return err
		}

		cfg.ReportInterval = reportInterval
		return nil
	})

	flag.Func("p", "Poll Interval", func(pollIntervalPlain string) error {
		if _, ok := os.LookupEnv("POLL_INTERVAL"); ok {
			return nil
		}

		pollInterval, err := time.ParseDuration(pollIntervalPlain)
		if err != nil {
			return err
		}

		cfg.PollInterval = pollInterval
		return nil
	})

	flag.Parse()
}
