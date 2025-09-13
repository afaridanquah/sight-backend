package customerapi

import (
	"bitbucket.org/msafaridanquah/sight-backend/app/domain/customerapp"
	"bitbucket.org/msafaridanquah/sight-backend/business/domain/customerbus"
	"bitbucket.org/msafaridanquah/sight-backend/business/domain/customerbus/postgres"
	"bitbucket.org/msafaridanquah/sight-backend/business/domain/documentbus"
	"bitbucket.org/msafaridanquah/sight-backend/foundation/envvar"
	"bitbucket.org/msafaridanquah/sight-backend/foundation/logger"
	"bitbucket.org/msafaridanquah/sight-backend/foundation/vaulti"
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
