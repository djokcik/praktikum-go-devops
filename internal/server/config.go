package server

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"github.com/djokcik/praktikum-go-devops/pkg/logging"
	"time"
)

type Config struct {
	Address       string        `env:"ADDRESS"`
	StoreInterval time.Duration `env:"STORE_INTERVAL"`
	StoreFile     string        `env:"STORE_FILE"`
	Restore       bool          `env:"RESTORE"`
}

func NewConfig() *Config {
	cfg := Config{
		Address:       "127.0.0.1:8080",
		StoreInterval: 300 * time.Second,
		Restore:       true,
		StoreFile:     "/tmp/devops-metrics-db.json",
	}

	cfg.parseFlags()
	cfg.parseEnv()

	return &cfg
}

func (cfg *Config) parseEnv() {
	err := env.Parse(cfg)
	if err != nil {
		logging.NewLogger().Fatal().Err(err).Msg("error parse environment")
	}
}

func (cfg *Config) parseFlags() {
	flag.StringVar(&cfg.Address, "a", cfg.Address, "Server address")
	flag.DurationVar(&cfg.StoreInterval, "i", cfg.StoreInterval, "Store save interval")
	flag.StringVar(&cfg.StoreFile, "f", cfg.StoreFile, "Store file")
	flag.BoolVar(&cfg.Restore, "r", cfg.Restore, "Restore")

	flag.Parse()
}
