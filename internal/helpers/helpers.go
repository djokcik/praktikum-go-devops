package helpers

import (
	"time"
)

func SetTicker(fn func(), seconds time.Duration) {
	ticker := time.NewTicker(seconds)

	for range ticker.C {
		fn()
	}
}

func Filter(arr []interface{}, f func(interface{}) bool) []interface{} {
	filtered := make([]interface{}, 0)
	for _, v := range arr {
		if f(v) {
			filtered = append(filtered, v)
		}
	}

	return filtered
}
