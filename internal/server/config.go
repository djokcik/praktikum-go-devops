package server

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"flag"
	"github.com/caarlos0/env/v6"
	"github.com/djokcik/praktikum-go-devops/pkg/logging"
	"net"
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
	GRPCAddress   string        `env:"GRPC_ADDRESS"`
	TrustedSubnet *net.IPNet
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
	var err error

	if mask := os.Getenv("TRUSTED_SUBNET"); mask != "" {
		cfg.TrustedSubnet = parseCIDR(mask)
	}

	if path := os.Getenv("CONFIG"); path != "" {
		err = cfg.parseConfigFile(path)
		if err != nil {
			logging.NewLogger().Fatal().Err(err).Msg("error parse config file")
		}
	}

	err = env.Parse(cfg)
	if err != nil {
		logging.NewLogger().Fatal().Err(err).Msg("error parse environment")
	}
}

func (cfg *Config) parseFlags() {
	flag.Func("c", "config path", func(s string) error {
		if os.Getenv("CONFIG") != "" {
			return nil
		}

		return cfg.parseConfigFile(s)
	})

	flag.StringVar(&cfg.Address, "a", cfg.Address, "Server address")
	flag.StringVar(&cfg.Key, "k", cfg.Key, "Hash key")
	flag.StringVar(&cfg.DatabaseDsn, "d", cfg.Key, "Database dsn")
	flag.DurationVar(&cfg.StoreInterval, "i", cfg.StoreInterval, "Store save interval")
	flag.StringVar(&cfg.StoreFile, "f", cfg.StoreFile, "Store file")
	flag.BoolVar(&cfg.Restore, "r", cfg.Restore, "Restore")
	flag.StringVar(&cfg.GRPCAddress, "g", cfg.GRPCAddress, "gRPC server address")
	flag.Func("t", "CIDR", func(mask string) error {
		cfg.TrustedSubnet = parseCIDR(mask)

		return nil
	})

	flag.Func("pk", "Private key file", func(s string) error {
		priv, err := parseRsaPrivateKeyFromPemStr(s)
		if err != nil {
			return err
		}

		cfg.PrivateKey.PrivateKey = priv
		return nil
	})

	flag.Parse()
}

type serverConfigFile struct {
	Address       string `json:"address"`
	Restore       bool   `json:"restore"`
	StoreInterval string `json:"store_interval"`
	StoreFile     string `json:"store_file"`
	DatabaseDsn   string `json:"database_dsn"`
	CryptoKey     string `json:"crypto_key"`
	TrustedSubnet string `json:"trusted_subnet"`
	GRPCAddress   string `json:"grpc_address"`
}

func (cfg *Config) parseConfigFile(path string) error {
	open, err := os.Open(path)
	if err != nil {
		return err
	}

	var config serverConfigFile
	err = json.NewDecoder(open).Decode(&config)
	if err != nil {
		return err
	}

	cfg.Address = config.Address
	cfg.Restore = config.Restore
	cfg.StoreFile = config.StoreFile
	cfg.DatabaseDsn = config.DatabaseDsn
	cfg.GRPCAddress = config.GRPCAddress
	cfg.TrustedSubnet = parseCIDR(config.TrustedSubnet)
	cfg.StoreInterval, err = time.ParseDuration(config.StoreInterval)
	if err != nil {
		return err
	}

	pub, err := parseRsaPrivateKeyFromPemStr(config.CryptoKey)
	if err != nil {
		return err
	}

	cfg.PrivateKey.PrivateKey = pub

	return nil
}

func parseCIDR(mask string) *net.IPNet {
	if mask == "" {
		return nil
	}

	_, ipnet, err := net.ParseCIDR(mask)
	if err != nil {
		logging.NewLogger().Fatal().Err(err).Msgf("invalid parse cidr: %s", mask)
	}

	return ipnet
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
