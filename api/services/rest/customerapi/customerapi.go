package customerapi

import (
	"bitbucket.org/msafaridanquah/sight-backend/app/domain/customerapp"
	"bitbucket.org/msafaridanquah/sight-backend/business/domain/customerbus"
	"bitbucket.org/msafaridanquah/sight-backend/business/domain/customerbus/postgres"
	"bitbucket.org/msafaridanquah/sight-backend/business/domain/documentbus"
	dr "bitbucket.org/msafaridanquah/sight-backend/business/domain/documentbus/postgres"
	"bitbucket.org/msafaridanquah/sight-backend/business/domain/identificationbus"
	ir "bitbucket.org/msafaridanquah/sight-backend/business/domain/identificationbus/postgres"
	"bitbucket.org/msafaridanquah/sight-backend/foundation/envvar"
	"bitbucket.org/msafaridanquah/sight-backend/foundation/logger"
	"bitbucket.org/msafaridanquah/sight-backend/foundation/vaulti"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nats-io/nats.go"
)

type CustomerApiConfig struct {
	Log      *logger.Logger
	Pool     *pgxpool.Pool
	Envvar   *envvar.Configuration
	Vault    *vaulti.Vaulty
	Mux      chi.Router
	NatsConn *nats.Conn
}

func Routes(cfg CustomerApiConfig) {
	repo := postgres.New(cfg.Pool, cfg.Pool, cfg.Vault)
	drepo := dr.New(cfg.Pool, cfg.Pool, cfg.Vault)
	idrepo := ir.New(cfg.Pool, cfg.Pool, cfg.Vault)
	// msg := natsio.New(cfg.NatsConn, cfg.Log)

	service, _ := customerbus.New(repo, cfg.Log)
	ds, _ := documentbus.New(drepo, cfg.Log)
	ids, _ := identificationbus.New(idrepo, cfg.Log)

	customerapp.Register(customerapp.Config{
		Log:                   cfg.Log,
		Service:               service,
		DocumentService:       ds,
		Router:                cfg.Mux,
		Repo:                  repo,
		EnvVar:                cfg.Envvar,
		IdentificationService: ids,
	})
}
