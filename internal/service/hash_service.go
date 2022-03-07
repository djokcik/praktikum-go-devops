// Package service provides common services both client and server
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

// HashService interface which provide operations for metrics hash
type HashService interface {
	GetHash(ctx context.Context, str string) string
	GetCounterHash(ctx context.Context, name string, value metric.Counter) string
	GetGaugeHash(ctx context.Context, name string, value metric.Gauge) string
	Verify(ctx context.Context, expectedHash string, actualHash string) bool
}

// NewHashService constructor for HashService with hashKey
// if hashKey is empty string it means hashKey is not use
func NewHashService(hashKey string) HashService {
	return &hashServiceImpl{HashKey: hashKey}
}

type hashServiceImpl struct {
	HashKey string
}

// GetHash return hash from any string
func (s hashServiceImpl) GetHash(_ context.Context, str string) string {
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

// GetCounterHash return counter metric hash from name and value
func (s hashServiceImpl) GetCounterHash(ctx context.Context, name string, value metric.Counter) string {
	return s.GetHash(ctx, fmt.Sprintf("%s:counter:%d", name, value))
}

// GetGaugeHash return gauge metric hash from name and value
func (s hashServiceImpl) GetGaugeHash(ctx context.Context, name string, value metric.Gauge) string {
	return s.GetHash(ctx, fmt.Sprintf("%s:gauge:%f", name, value))
}

// Verify checks equals expected and actual hash
func (s hashServiceImpl) Verify(ctx context.Context, expectedHash string, actualHash string) bool {
	if s.HashKey == "" {
		return true
	}

	if !hmac.Equal([]byte(expectedHash), []byte(actualHash)) {
		s.Log(ctx).Warn().Msgf("Invalid verified hash: expected: %s, actual: %s", expectedHash, actualHash)
		return false
	}

	return true
}

// Log return zerolog with configure service key
func (s hashServiceImpl) Log(ctx context.Context) *zerolog.Logger {
	_, logger := logging.GetCtxLogger(ctx)
	logger = logger.With().Str(logging.ServiceKey, "hash service").Logger()

	return &logger
}
