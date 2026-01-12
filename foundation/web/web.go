package web

import (
	"context"
	"net/http"

	"go.opentelemetry.io/otel/trace"
)

type Logger func(ctx context.Context, msg string, args ...any)

type MidFunc func(next http.Handler) http.Handler

type App struct {
	Address string
	// PostgresDB  *pgxpool.Pool
	// Metrics     http.Handler
	Mw     []MidFunc
	Log    Logger
	Tracer trace.Tracer
	// Mux    *chi.Mux

	// Vault       vaulti.Vaulty
}

// NewApp creates an App value that handle a set of routes for the application.
func NewApp(log Logger, tracer trace.Tracer, mw ...MidFunc) *App {
	// Create an OpenTelemetry HTTP Handler which wraps our router. This will start
	// the initial span and annotate it with information about the request/trusted.
	//
	// This is configured to use the W3C TraceContext standard to set the remote
	// parent if a client request includes the appropriate headers.
	// https://w3c.github.io/trace-context/

	// mux := chi.NewMux()

	return &App{
		Log:    log,
		Tracer: tracer,
		// Mux:    mux,
		Mw: mw,
	}
}
