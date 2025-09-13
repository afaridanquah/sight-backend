package businessapi

import (
	"bitbucket.org/msafaridanquah/verifylab-service/app/domain/businessapp"
	"bitbucket.org/msafaridanquah/verifylab-service/business/domain/businessbus"
	"bitbucket.org/msafaridanquah/verifylab-service/business/domain/businessbus/postgres"
	"bitbucket.org/msafaridanquah/verifylab-service/foundation/envvar"
	"bitbucket.org/msafaridanquah/verifylab-service/foundation/logger"
	"bitbucket.org/msafaridanquah/verifylab-service/foundation/vaulti"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Routes(logger *logger.Logger, pool *pgxpool.Pool, envvar *envvar.Configuration, vaulti *vaulti.Vaulty, chi chi.Router) {
	repo := postgres.New(pool, pool, vaulti)

	service, _ := businessbus.New(repo, logger)

	businessapp.Register(businessapp.Config{
		Log:     logger,
		Service: service,
		Router:  chi,
		Repo:    repo,
	})
}
