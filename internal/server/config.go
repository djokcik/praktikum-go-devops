package server

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"flag"
	"github.com/caarlos0/env/v6"
	"github.com/djokcik/praktikum-go-devops/pkg/logging"
	"os"
	"time"
)

type RsaPrivateKey struct {
	*rsa.PrivateKey
}

func (r *RsaPrivateKey) UnmarshalText(text []byte) error {
	priv, err := parseRsaPrivateKeyFromPemStr(string(text))
	if err != nil {
		return err
	}

	r.PrivateKey = priv
	return nil
}

type Config struct {
	Address       string        `env:"ADDRESS"`
	StoreInterval time.Duration `env:"STORE_INTERVAL"`
	StoreFile     string        `env:"STORE_FILE"`
	Restore       bool          `env:"RESTORE"`
	Key           string        `env:"KEY"`
	DatabaseDsn   string        `env:"DATABASE_DSN"`
	PrivateKey    RsaPrivateKey `env:"CRYPTO_KEY"`
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
	flag.DurationVar(&cfg.StoreInterval, "i", cfg.StoreInterval, "Store save interval")
	flag.StringVar(&cfg.StoreFile, "f", cfg.StoreFile, "Store file")
	flag.BoolVar(&cfg.Restore, "r", cfg.Restore, "Restore")
	flag.Func("c", "Private key file", func(s string) error {
		priv, err := parseRsaPrivateKeyFromPemStr(s)
		if err != nil {
			return err
		}

		cfg.PrivateKey.PrivateKey = priv
		return nil
	})

	flag.Parse()
}

func parseRsaPrivateKeyFromPemStr(fileName string) (*rsa.PrivateKey, error) {
	privPEM, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(privPEM)
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}

	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return priv, nil
}
