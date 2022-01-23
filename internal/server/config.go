package server

import (
	"flag"
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/djokcik/praktikum-go-devops/pkg/logging"
	"time"
)

type Config struct {
	Address       string        `env:"ADDRESS"`
	StoreInterval time.Duration `env:"STORE_INTERVAL"`
	StoreFile     string        `env:"STORE_FILE"`
	Restore       bool          `env:"RESTORE"`
	Key           string        `env:"KEY"`
	DatabaseDsn   string        `env:"DATABASE_DSN"`
}

func NewConfig() Config {
	cfg := Config{
		Address:       "127.0.0.1:8080",
		StoreInterval: 300 * time.Second,
		Restore:       true,
		StoreFile:     "/tmp/devops-metrics-db.json",
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
	flag.StringVar(&cfg.DatabaseDsn, "d", cfg.Key, "Database dsn")
	flag.DurationVar(&cfg.StoreInterval, "i", cfg.StoreInterval, "store save interval")
	flag.StringVar(&cfg.StoreFile, "f", cfg.StoreFile, "store file")
	flag.BoolVar(&cfg.Restore, "r", cfg.Restore, "Restore")

	flag.Parse()
}

func (cfg Config) String() string {
	return fmt.Sprintf("Start Server. "+
		"Address: %s, "+
		"StoreInterval: %s, "+
		"StoreFile: %s, "+
		"Restore: %v, "+
		"DatabaseDsn: %s",
		cfg.Address, cfg.StoreInterval, cfg.StoreFile, cfg.Restore, cfg.DatabaseDsn)
}
