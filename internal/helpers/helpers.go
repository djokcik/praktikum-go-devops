// Package helpers provides any helper functions
package helpers

import (
	"time"
)

// SetTicker each `seconds` call fn
func SetTicker(fn func(), seconds time.Duration) {
	ticker := time.NewTicker(seconds)

	for range ticker.C {
		fn()
	}
}
