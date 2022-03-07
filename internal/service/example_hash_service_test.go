package service_test

import (
	"context"
	"fmt"
	"github.com/djokcik/praktikum-go-devops/internal/metric"
	"github.com/djokcik/praktikum-go-devops/internal/service"
	"github.com/rs/zerolog"
	"time"
)

var _ = 1

func Example() {
	zerolog.SetGlobalLevel(zerolog.Disabled)

	service := service.NewHashService("MyHashKey")
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	hash := service.GetHash(ctx, "cpu:counter:10")
	fmt.Printf("GetHash: %s\n", hash)

	counterHash := service.GetCounterHash(ctx, "cpu", metric.Counter(10))
	fmt.Printf("GetCounterHash: %s\n", counterHash)

	gaugeHash := service.GetGaugeHash(ctx, "cpu", metric.Gauge(10))
	fmt.Printf("GetCounterHash: %s\n", gaugeHash)

	fmt.Printf("Verify CounterHash and GaugeHash: %v\n", service.Verify(ctx, counterHash, gaugeHash))

	// Output: GetHash: 558eb947e5c8cec059797cc8966f751cb74aa0e75af3ab6b3fc07d099d134141
	// GetCounterHash: 558eb947e5c8cec059797cc8966f751cb74aa0e75af3ab6b3fc07d099d134141
	// GetCounterHash: 20bf0e46f7956d976fedb1811422cba1206bf13bc7d594004e93ce290bd93bcd
	// Verify CounterHash and GaugeHash: false
}
