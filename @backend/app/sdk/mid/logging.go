package mid

import (
	"log/slog"
	"net/http"
	"time"

	"bitbucket.org/msafaridanquah/verifylab-service/foundation/logger"
)

func Logger(logger *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			now := time.Now()
			logger.Info(r.Context(), "request started",
				slog.Time("started_time", now),
				slog.String("method", r.Method),
				slog.String("url", r.URL.String()),
				slog.String("remoteaddr", r.RemoteAddr),
			)

			logger.Info(r.Context(), "request completed",
				slog.String("method", r.Method),
				slog.String("remoteaddr", r.RemoteAddr),
				slog.String("since", time.Since(now).String()),
			)

			next.ServeHTTP(w, r)
		})
	}
}
