package helpers

import "time"

type ContextKey string

const (
	ShutdownServerSyncContext ContextKey = "ShutdownServerSyncContext"
)

func SetTicker(fn func(), seconds time.Duration) {
	ticker := time.NewTicker(seconds)

	for range ticker.C {
		fn()
	}
}
