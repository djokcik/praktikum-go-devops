package service

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/djokcik/praktikum-go-devops/internal/metric"
	"github.com/djokcik/praktikum-go-devops/pkg/logging"
	"github.com/rs/zerolog"
)

//go:generate mockery --name=HashService

type HashService interface {
	GetHash(ctx context.Context, str string) string
	GetCounterHash(ctx context.Context, name string, value metric.Counter) string
	GetGaugeHash(ctx context.Context, name string, value metric.Gauge) string
	Verify(ctx context.Context, expectedHash string, actualHash string) bool
}

type HashServiceImpl struct {
	HashKey string
}

func (s HashServiceImpl) GetHash(ctx context.Context, str string) string {
	if s.HashKey == "" {
		return ""
	}

	h := hmac.New(sha256.New, []byte(s.HashKey))
	h.Write([]byte(str))
	sign := h.Sum(nil)

	hash := make([]byte, hex.EncodedLen(len(sign)))
	hex.Encode(hash, sign)

	return string(hash)
}

func (s HashServiceImpl) GetCounterHash(ctx context.Context, name string, value metric.Counter) string {
	return s.GetHash(ctx, fmt.Sprintf("%s:counter:%d", name, value))
}

func (s HashServiceImpl) GetGaugeHash(ctx context.Context, name string, value metric.Gauge) string {
	return s.GetHash(ctx, fmt.Sprintf("%s:gauge:%f", name, value))
}

func (s HashServiceImpl) Verify(ctx context.Context, expectedHash string, actualHash string) bool {
	if s.HashKey == "" {
		return true
	}

	if !hmac.Equal([]byte(expectedHash), []byte(actualHash)) {
		s.Log(ctx).Warn().Msgf("Invalid verified hash: expected: %s, actual: %s", expectedHash, actualHash)
		return false
	}

	return true
}

func (s HashServiceImpl) Log(ctx context.Context) *zerolog.Logger {
	_, logger := logging.GetCtxLogger(ctx)
	logger = logger.With().Str(logging.ServiceKey, "hash service").Logger()

	return &logger
}
