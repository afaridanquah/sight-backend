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

	"bitbucket.org/msafaridanquah/sight-backend/app/sdk"
	"bitbucket.org/msafaridanquah/sight-backend/app/sdk/mid"
	"bitbucket.org/msafaridanquah/sight-backend/foundation/envvar"
	"bitbucket.org/msafaridanquah/sight-backend/foundation/logger"
	"bitbucket.org/msafaridanquah/sight-backend/foundation/otel"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const serviceName = "IMPORT_FROM_JSON"

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
	//
	cfg := struct {
		Address     string
		PostgresDB  *pgxpool.Pool
		Metrics     http.Handler
		Middlewares []func(next http.Handler) http.Handler
		Logger      *logger.Logger
		// Tracer      trace.Tracer
		// Vault       vaulti.Vaulty
	}{
		Address:     ":" + address,
		Metrics:     promhttp.Handler(),
		Middlewares: []func(next http.Handler) http.Handler{mid.Otel(tracer), mid.Logger(log)},
		Logger:      log,
		PostgresDB:  pool,
		// Tracer:      tracer,
	}

	files, err := os.ReadDir("/Users/afaridanquah/code/golang/sight/exports")
	if err != nil {
		return fmt.Errorf("failed to read directory: %w", err)
	}

	fmt.Printf("first file %+v", files[0])

	// -------------------------------------------------------------------------
	// Start API Service

	log.Info(ctx, "startup", "status", "initializing V1 API support")

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	router := chi.NewRouter()

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

type Identification struct {
	UID                int       `json:"uid"`
	Forenames          string    `json:"forenames"`
	Surname            string    `json:"surname"`
	PrevOrMaidenName   string    `json:"prev_or_maiden_name"`
	Sex                string    `json:"sex"`
	Occupation         string    `json:"occupation"`
	MaritalStatus      string    `json:"marital_status"`
	DateOfBirth        string    `json:"date_of_birth"`
	BirthTown          string    `json:"birth_town"`
	BirthCountry       string    `json:"birth_country"`
	BirthRegion        string    `json:"birth_region"`
	BirthDistrict      string    `json:"birth_district"`
	Nationality        string    `json:"nationality"`
	Resident           string    `json:"resident"`
	SocialSecurityNo   string    `json:"social_security_no"`
	MotherMaidenName   string    `json:"mother_maiden_name"`
	MotherForename     string    `json:"mother_forename"`
	NationalID         string    `json:"national_id"`
	IssueDate          string    `json:"issue_date"`
	ExpiryDate         string    `json:"expiry_date"`
	CountryOfIssue     string    `json:"country_of_issue"`
	PlaceOfIssue       string    `json:"place_of_issue"`
	CardNumber         string    `json:"card_number"`
	HouseNumber        string    `json:"house_number"`
	StreetName         string    `json:"street_name"`
	Town               string    `json:"town"`
	City               string    `json:"city"`
	Community          string    `json:"community"`
	Country            string    `json:"country"`
	Region             string    `json:"region"`
	District           string    `json:"district"`
	PostalAddress      string    `json:"postal_address"`
	PhoneNumber1       string    `json:"phone_number_1"`
	PhoneNumber2       any       `json:"phone_number_2"`
	Email              string    `json:"email"`
	TinNumber          string    `json:"tin_number"`
	TinIssuedDate      any       `json:"tin_issued_date"`
	LastUpdate         time.Time `json:"last_update"`
	CreatedDate        time.Time `json:"created_date"`
	DigLongitude       any       `json:"dig_Longitude"`
	DigLatitude        any       `json:"dig_Latitude"`
	DigStreet          any       `json:"dig_Street"`
	DigRegion          any       `json:"dig_Region"`
	DigArea            any       `json:"dig_Area"`
	DigDistrict        any       `json:"dig_District"`
	DigPostCode        any       `json:"dig_PostCode"`
	GhanaPostAddress   string    `json:"ghana_post_address"`
	LocationUID        any       `json:"location_uid"`
	RefisicID          any       `json:"refisicID"`
	RegistrationSTATUS string    `json:"registration_STATUS"`
}
