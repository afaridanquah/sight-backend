package web

import (
	"net/http"

	"bitbucket.org/msafaridanquah/verifylab-service/foundation/logger"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.opentelemetry.io/otel/trace"
)

type Config struct {
	Name        string
	Address     string
	PostgresDB  *pgxpool.Pool
	Metrics     http.Handler
	Middlewares []func(next http.Handler) http.Handler
	Logger      *logger.Logger
	Tracer      trace.Tracer
}
