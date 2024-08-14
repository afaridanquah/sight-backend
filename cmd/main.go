package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"bitbucket.org/msafaridanquah/verifylab-service/cmd/internal"
	internaldomain "bitbucket.org/msafaridanquah/verifylab-service/internal"
	"bitbucket.org/msafaridanquah/verifylab-service/internal/envvar"
	"bitbucket.org/msafaridanquah/verifylab-service/internal/rest"
	"bitbucket.org/msafaridanquah/verifylab-service/internal/service/customer"
	"bitbucket.org/msafaridanquah/verifylab-service/internal/service/verification"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/riandyrn/otelchi"
	"go.mongodb.org/mongo-driver/mongo"
)

func main() {
	var env, address string

	flag.StringVar(&env, "env", "", "Environment Variables filename")
	flag.StringVar(&address, "address", ":9235", "HTTP Server Address")
	flag.Parse()

	errC, err := run(env, address)
	if err != nil {
		log.Fatalf("Couldn't run: %s", err)
	}

	if err := <-errC; err != nil {
		log.Fatalf("Error while running: %s", err)
	}

}

type serverConfig struct {
	Address     string
	MongoDB     *mongo.Client
	Metrics     http.Handler
	Middlewares []func(next http.Handler) http.Handler
	Logger      *slog.Logger
}

func run(env, address string) (<-chan error, error) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	if err := envvar.Load(env); err != nil {
		return nil, internaldomain.WrapErrorf(err, internaldomain.ErrorCodeUnknown, "envvar.Load")
	}

	vault, err := internal.NewVaultProvider()
	if err != nil {
		return nil, internaldomain.WrapErrorf(err, internaldomain.ErrorCodeUnknown, "internal.NewVaultProvider")
	}

	conf := envvar.New(vault)

	_, err = internal.NewOTExporter(conf)

	if err != nil {
		return nil, internaldomain.WrapErrorf(err, internaldomain.ErrorCodeUnknown, "internal.NewOTExporter")
	}

	mongoDB, err := internal.NewMongoDB(conf)

	if err != nil {
		return nil, internaldomain.WrapErrorf(err, internaldomain.ErrorCodeUnknown, "internal.NewMongoDB")
	}

	logging := func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Info(r.Method,
				slog.Time("time", time.Now()),
				slog.String("url", r.URL.String()),
			)

			h.ServeHTTP(w, r)
		})
	}

	srv, err := newServer(serverConfig{
		Address:     ":" + address,
		MongoDB:     mongoDB,
		Metrics:     promhttp.Handler(),
		Middlewares: []func(next http.Handler) http.Handler{otelchi.Middleware("verifylab-api-server"), logging},
		Logger:      logger,
	})

	log.Printf("Server starting on port %v", address)

	if err != nil {
		return nil, internaldomain.WrapErrorf(err, internaldomain.ErrorCodeUnknown, "newServer")
	}

	errC := make(chan error, 1)

	ctx, stop := signal.NotifyContext(context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	go func() {
		<-ctx.Done()

		logger.Info("Shutdown signal received")

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

		logger.Info("Shutdown completed")
	}()

	go func() {
		logger.Info("Listening and serving", slog.String("address", address))

		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errC <- err
		}
	}()

	return errC, nil

}

func newServer(conf serverConfig) (*http.Server, error) {
	router := chi.NewRouter()
	router.Use(render.SetContentType(render.ContentTypeJSON))

	for _, mw := range conf.Middlewares {
		router.Use(mw)
	}

	v1Router := chi.NewRouter()

	v1Router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	cs, err := customer.New(conf.Logger, customer.WithMemoryRepository())
	if err != nil {
		return nil, internaldomain.WrapErrorf(err, internaldomain.ErrorCodeUnknown, "customer service")
	}
	vs, err := verification.New(conf.Logger, verification.WithMemoryVerificationRepository())

	if err != nil {
		return nil, internaldomain.WrapErrorf(err, internaldomain.ErrorCodeUnknown, "verification service")
	}

	rest.NewCustomerhandler(*cs).Register(v1Router)
	rest.NewVerificationHandler(*vs).Register(v1Router)

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
