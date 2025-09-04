package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"bitbucket.org/msafaridanquah/verifylab-service/api/cmd/domain/http/businessapi"
	"bitbucket.org/msafaridanquah/verifylab-service/api/cmd/domain/http/customerapi"
	"bitbucket.org/msafaridanquah/verifylab-service/api/cmd/domain/http/otpapi"
	"bitbucket.org/msafaridanquah/verifylab-service/api/cmd/domain/http/verificationapi"
	"bitbucket.org/msafaridanquah/verifylab-service/app/sdk"
	"bitbucket.org/msafaridanquah/verifylab-service/app/sdk/mid"
	"bitbucket.org/msafaridanquah/verifylab-service/foundation/envvar"
	"bitbucket.org/msafaridanquah/verifylab-service/foundation/logger"
	"bitbucket.org/msafaridanquah/verifylab-service/foundation/otel"
	"bitbucket.org/msafaridanquah/verifylab-service/foundation/vaulti"
	"bitbucket.org/msafaridanquah/verifylab-service/foundation/web"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/riandyrn/otelchi"
)

const serviceName = "SIGHT"

func main() {
	var log *logger.Logger
	var env, address string

	flag.StringVar(&env, "env", "env.example", "Environment Variables filename")
	flag.StringVar(&address, "address", "9235", "HTTP Server Address")
	flag.Parse()

	ctx := context.Background()

	traceIDFn := func(ctx context.Context) string {
		return otel.GetTraceID(ctx)
	}

	log = logger.New(os.Stdout, logger.LevelInfo, serviceName, traceIDFn)

	// -------------------------------------------------------------------------

	if err := run(ctx, env, address, log); err != nil {
		log.Error(ctx, "startup", "err", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, env string, address string, log *logger.Logger) error {
	service := serviceName

	// -------------------------------------------------------------------------
	// GOMAXPROCS
	log.Info(ctx, "startup", "GOMAXPROCS", runtime.GOMAXPROCS(0))
	log.Info(ctx, "starting service")
	defer log.Info(ctx, "shutdown complete")

	if err := envvar.Load(env); err != nil {
		return fmt.Errorf("load envvar %w", err)
	}

	vault, err := sdk.NewVaultProvider()
	if err != nil {
		return fmt.Errorf("new vault provider %w", err)
	}

	conf := envvar.New(vault)

	var tempo otel.Config

	jaegerEndpoint, _ := conf.Get("JAEGER_ENDPOINT")

	tempo = otel.Config{
		Host:        jaegerEndpoint,
		ServiceName: service,
		Probability: 0.05,
		ExcludedRoutes: map[string]struct{}{
			"/v1/liveness":  {},
			"/v1/readiness": {},
		},
	}

	traceProvider, teardown, err := otel.InitTracing(log, tempo)
	if err != nil {
		return fmt.Errorf("init tracing: %w", err)
	}

	defer teardown(ctx)
	tracer := traceProvider.Tracer(tempo.ServiceName)

	pool, err := sdk.NewPostgreSQL(ctx, conf)
	if err != nil {
		return fmt.Errorf("new postgres sql %w", err)
	}

	// -------------------------------------------------------------------------
	// Configuration
	cfg := web.Config{
		Name:        "",
		Address:     ":" + address,
		Metrics:     promhttp.Handler(),
		Middlewares: []func(next http.Handler) http.Handler{mid.Otel(tracer), mid.Logger(log)},
		Logger:      log,
		PostgresDB:  pool,
		Tracer:      tracer,
	}

	// -------------------------------------------------------------------------
	// Start API Service

	log.Info(ctx, "startup", "status", "initializing V1 API support")

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	router := chi.NewRouter()
	router.Use(render.SetContentType(render.ContentTypeJSON))
	router.Use(otelchi.Middleware(tempo.ServiceName, otelchi.WithChiRoutes(router)))

	for _, mw := range cfg.Middlewares {
		router.Use(mw)
	}

	v1Router := chi.NewRouter()

	v1Router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	vaulti, err := vaulti.InitVault(vaulti.Config{
		Log: cfg.Logger,
	})

	if err != nil {
		return err
	}

	businessapi.Routes(cfg.Logger, cfg.PostgresDB, conf, vaulti, v1Router)
	customerapi.Routes(cfg.Logger, cfg.PostgresDB, conf, vaulti, v1Router)
	verificationapi.Routes(cfg.Logger, cfg.PostgresDB, conf, vaulti, v1Router)
	otpapi.Routes(cfg.Logger, cfg.PostgresDB, conf, vaulti, v1Router)

	router.Mount("/v1", v1Router)

	api := http.Server{
		Handler:           router,
		Addr:              cfg.Address,
		ReadTimeout:       1 * time.Second,
		ReadHeaderTimeout: 1 * time.Second,
		WriteTimeout:      1 * time.Second,
		IdleTimeout:       1 * time.Second,
	}

	serverErrors := make(chan error, 1)

	go func() {
		log.Info(ctx, "startup", "status", "api router started", "host", api.Addr)
		serverErrors <- api.ListenAndServe()
	}()

	// -------------------------------------------------------------------------
	// Shutdown

	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)

	case sig := <-shutdown:
		log.Info(ctx, "shutdown", "status", "shutdown started", "signal", sig)
		defer log.Info(ctx, "shutdown", "status", "shutdown complete", "signal", sig)

		ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
		defer cancel()

		if err := api.Shutdown(ctx); err != nil {
			defer func() {
				if err := api.Close(); err != nil {
					return
				}
			}()
			return fmt.Errorf("could not stop server gracefully: %w", err)
		}
	}

	return nil

}
