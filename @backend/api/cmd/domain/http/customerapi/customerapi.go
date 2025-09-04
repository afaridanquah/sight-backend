package customerapi

import (
	"bitbucket.org/msafaridanquah/verifylab-service/app/domain/customerapp"
	"bitbucket.org/msafaridanquah/verifylab-service/business/domain/customerbus"
	"bitbucket.org/msafaridanquah/verifylab-service/business/domain/customerbus/postgres"
	"bitbucket.org/msafaridanquah/verifylab-service/business/domain/documentbus"
	"bitbucket.org/msafaridanquah/verifylab-service/foundation/envvar"
	"bitbucket.org/msafaridanquah/verifylab-service/foundation/logger"
	"bitbucket.org/msafaridanquah/verifylab-service/foundation/vaulti"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Routes(logger *logger.Logger, pool *pgxpool.Pool, envvar *envvar.Configuration, vaulti *vaulti.Vaulty, chi chi.Router) {
	repo := postgres.New(pool, pool, vaulti)

	service, _ := customerbus.New(repo, logger)
	ds, _ := documentbus.New(logger)

	customerapp.Register(customerapp.Config{
		Log:             logger,
		Service:         service,
		DocumentService: ds,
		Router:          chi,
		Repo:            repo,
		EnvVar:          envvar,
	})
}
