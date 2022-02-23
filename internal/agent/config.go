package agent

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"github.com/djokcik/praktikum-go-devops/pkg/logging"
	"time"
)

type Config struct {
	Address             string        `env:"ADDRESS"`
	ReportInterval      time.Duration `env:"REPORT_INTERVAL"`
	PollInterval        time.Duration `env:"POLL_INTERVAL"`
	PollPsutilsInterval time.Duration `envDefault:"10s"`
	Key                 string        `env:"KEY"`
}

func NewConfig() Config {
	cfg := Config{
		Address:        "127.0.0.1:8080",
		ReportInterval: 10 * time.Second,
		PollInterval:   2 * time.Second,
	}

	cfg.parseFlags()
	cfg.parseEnv()

	return cfg
}

func (cfg *Config) parseEnv() {
	err := env.Parse(cfg)
	if err != nil {
		logging.NewLogger().Fatal().Err(err).Msg("error parse environment")
	}
}

func (cfg *Config) parseFlags() {
	flag.StringVar(&cfg.Address, "a", cfg.Address, "Server address")
	flag.StringVar(&cfg.Key, "k", cfg.Key, "Hash key")
	flag.DurationVar(&cfg.ReportInterval, "r", cfg.ReportInterval, "Report Interval")
	flag.DurationVar(&cfg.PollInterval, "p", cfg.PollInterval, "Poll Interval")

	flag.Parse()
}
