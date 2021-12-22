package helpers

import "time"

func SetTicker(fn func(), seconds time.Duration) {
	ticker := time.NewTicker(seconds * time.Second)

	for range ticker.C {
		fn()
	}
}
