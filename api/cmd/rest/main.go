package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"bitbucket.org/msafaridanquah/verifylab-service/api/cmd/domain/http/customerapi"
	"bitbucket.org/msafaridanquah/verifylab-service/api/cmd/domain/http/identificationapi"
	"bitbucket.org/msafaridanquah/verifylab-service/api/cmd/domain/http/otpapi"
	"bitbucket.org/msafaridanquah/verifylab-service/app/sdk"
	"bitbucket.org/msafaridanquah/verifylab-service/app/sdk/mid"
	"bitbucket.org/msafaridanquah/verifylab-service/foundation/envvar"
	"bitbucket.org/msafaridanquah/verifylab-service/foundation/ierr"
	"bitbucket.org/msafaridanquah/verifylab-service/foundation/logger"
	"bitbucket.org/msafaridanquah/verifylab-service/foundation/otel"
	"bitbucket.org/msafaridanquah/verifylab-service/foundation/web"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	var env, address string

	flag.StringVar(&env, "env", "env.example", "Environment Variables filename")
	flag.StringVar(&address, "address", "9235", "HTTP Server Address")
	flag.Parse()

	ctx := context.Background()

	traceIDFn := func(ctx context.Context) string {
		return otel.GetTraceID(ctx)
	}

	log := logger.New(os.Stdout, logger.LevelInfo, "verifylab", traceIDFn)

	log.Info(ctx, "server starting on port", "address", address)

	errC, err := run(env, address, ctx, log)
	if err != nil {
		fmt.Printf("err: %v", err)
		os.Exit(1)
	}

	if err := <-errC; err != nil {
		fmt.Printf("Error while running: %s", err)
	}
}

func run(env, address string, ctx context.Context, log *logger.Logger) (<-chan error, error) {
	if err := envvar.Load(env); err != nil {

		return nil, fmt.Errorf("load envvar %w", err)
	}

	vault, err := sdk.NewVaultProvider()
	if err != nil {
		return nil, fmt.Errorf("new vault provider %w", err)
	}

	conf := envvar.New(vault)

	pool, err := sdk.NewPostgreSQL(conf)
	if err != nil {
		return nil, fmt.Errorf("new postgres sql %w", err)
	}

	var tempo otel.Config

	jaegerEndpoint, _ := conf.Get("JAEGER_ENDPOINT")

	tempo = otel.Config{
		Host:        jaegerEndpoint,
		ServiceName: "verifylab",
		Probability: 0.05,
		ExcludedRoutes: map[string]struct{}{
			"/v1/liveness":  {},
			"/v1/readiness": {},
		},
	}

	traceProvider, _, err := otel.InitTracing(log, tempo)
	tracer := traceProvider.Tracer(tempo.ServiceName)

	if err != nil {
		return nil, ierr.WrapErrorf(err, ierr.ErrorCodeUnknown, "sdk.NewOTExporter")
	}

	// defer teardown(context.Background())

	srv, err := newServer(web.Config{
		Name:        tempo.ServiceName,
		Address:     ":" + address,
		Metrics:     promhttp.Handler(),
		Middlewares: []func(next http.Handler) http.Handler{mid.Otel(tracer), mid.Logger(log)},
		Logger:      log,
		PostgresDB:  pool,
		Tracer:      tracer,
	})

	fmt.Printf("Server starting on port: %v", address)

	if err != nil {
		return nil, ierr.WrapErrorf(err, ierr.ErrorCodeUnknown, "newServer")
	}

	errC := make(chan error, 1)

	ctx, stop := signal.NotifyContext(context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	go func() {
		<-ctx.Done()

		log.Info(ctx, "Shutdown signal received")

		ctxTimeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)

		defer func() {
			stop()
			cancel()
			close(errC)
		}()

		srv.SetKeepAlivesEnabled(false)

		if err := srv.Shutdown(ctxTimeout); err != nil {
			errC <- err
		}

		log.Info(ctx, "Shutdown completed")
	}()

	go func() {
		log.Info(ctx, "Listening and serving", slog.String("address", address))

		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errC <- err
		}
	}()

	return errC, nil

}

func newServer(conf web.Config) (*http.Server, error) {
	router := chi.NewRouter()
	router.Use(render.SetContentType(render.ContentTypeJSON))

	for _, mw := range conf.Middlewares {
		router.Use(mw)
	}

	v1Router := chi.NewRouter()

	v1Router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	customerapi.Routes(conf.Logger, conf.PostgresDB, v1Router)
	identificationapi.Routes(conf.Logger, conf.PostgresDB, v1Router)
	otpapi.Routes(conf.Logger, conf.PostgresDB, v1Router)

	router.Mount("/v1", v1Router)
	router.Handle("/metrics", conf.Metrics)

	return &http.Server{
		Handler:           router,
		Addr:              conf.Address,
		ReadTimeout:       1 * time.Second,
		ReadHeaderTimeout: 1 * time.Second,
		WriteTimeout:      1 * time.Second,
		IdleTimeout:       1 * time.Second,
	}, nil
}
