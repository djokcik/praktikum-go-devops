package middleware

import (
	"github.com/djokcik/praktikum-go-devops/internal/server"
	"github.com/djokcik/praktikum-go-devops/pkg/logging"
	"net"
	"net/http"
)

func TrustedSubnetHandle(cfg server.Config) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if cfg.TrustedSubnet == nil {
				next.ServeHTTP(w, r)
				return
			}

			if r.RemoteAddr == "" {
				logging.NewLogger().Info().Msg("invalid real IP")
				http.Error(w, "invalid real IP", http.StatusBadRequest)
				return
			}

			realIP := net.ParseIP(r.RemoteAddr)
			if !cfg.TrustedSubnet.Contains(realIP) {
				logging.NewLogger().Info().Msg("cannot trusted IP")
				http.Error(w, "cannot trusted IP", http.StatusBadRequest)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
