package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"runtime"

	"bitbucket.org/msafaridanquah/verifylab-service/api/cmd/cli/command"
	"bitbucket.org/msafaridanquah/verifylab-service/app/sdk"
	"bitbucket.org/msafaridanquah/verifylab-service/foundation/envvar"
	"bitbucket.org/msafaridanquah/verifylab-service/foundation/logger"
	"bitbucket.org/msafaridanquah/verifylab-service/foundation/otel"
)

const serviceName = "PORTAL"

func main() {
	var log *logger.Logger
	var env string

	flag.StringVar(&env, "env", "env.example", "Environment Variables filename")
	// flag.StringVar(&address, "address", "9111", "HTTP Server Address")
	// flag.Parse()

	ctx := context.Background()

	traceIDFn := func(ctx context.Context) string {
		return otel.GetTraceID(ctx)
	}

	log = logger.New(os.Stdout, logger.LevelInfo, serviceName, traceIDFn)

	// -------------------------------------------------------------------------

	if err := run(ctx, env, log); err != nil {
		log.Error(ctx, "startup", "err", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, env string, log *logger.Logger) error {
	// service := serviceName

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
		ServiceName: serviceName,
		Probability: 0.05,
		ExcludedRoutes: map[string]struct{}{
			"/v1/liveness":  {},
			"/v1/readiness": {},
		},
	}

	traceProvider, _, err := otel.InitTracing(log, tempo)
	tracer := traceProvider.Tracer(tempo.ServiceName)
	if err != nil {
		return fmt.Errorf("init tracing: %w", err)
	}

	pool, err := sdk.NewPostgreSQL(conf)
	if err != nil {
		return fmt.Errorf("new postgres sql %w", err)
	}

	config := command.Config{
		Name:       "test",
		PostgresDB: pool,
		Tracer:     tracer,
	}

	command.New(config)
	return nil
}
