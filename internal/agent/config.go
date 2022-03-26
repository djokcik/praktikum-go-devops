package agent

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

type RsaPublicKey struct {
	*rsa.PublicKey
}

func (r *RsaPublicKey) UnmarshalText(text []byte) error {
	pub, err := parseRsaPublicKeyFromPemStr(string(text))
	if err != nil {
		return err
	}

	r.PublicKey = pub
	return nil
}

type Config struct {
	Address             string        `env:"ADDRESS"`
	ReportInterval      time.Duration `env:"REPORT_INTERVAL"`
	PollInterval        time.Duration `env:"POLL_INTERVAL"`
	PollPsutilsInterval time.Duration `envDefault:"10s"`
	Key                 string        `env:"KEY"`
	PublicKey           RsaPublicKey  `env:"CRYPTO_KEY"`
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
	flag.Func("c", "Public key file", func(s string) error {
		pub, err := parseRsaPublicKeyFromPemStr(s)
		if err != nil {
			return err
		}

		cfg.PublicKey.PublicKey = pub
		return nil
	})

	flag.Parse()
}

func parseRsaPublicKeyFromPemStr(text string) (*rsa.PublicKey, error) {
	pubPEM, err := os.ReadFile(text)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(pubPEM)
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	switch pub := pub.(type) {
	case *rsa.PublicKey:
		return pub, nil
	default:
		break // fall through
	}

	return nil, errors.New("key type is not RSA")
}
