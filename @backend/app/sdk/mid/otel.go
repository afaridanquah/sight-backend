package mid

import (
	"net/http"

	"bitbucket.org/msafaridanquah/verifylab-service/foundation/otel"
	"go.opentelemetry.io/otel/trace"
)

func Otel(tracer trace.Tracer) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := otel.InjectTracing(r.Context(), tracer)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
