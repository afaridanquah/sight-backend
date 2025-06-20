package identificationapi

import (
	"bitbucket.org/msafaridanquah/verifylab-service/app/domain/identificationapp"
	"bitbucket.org/msafaridanquah/verifylab-service/business/domain/identificationbus"
	"bitbucket.org/msafaridanquah/verifylab-service/business/domain/identificationbus/postgres"
	"bitbucket.org/msafaridanquah/verifylab-service/foundation/logger"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Routes(logger *logger.Logger, pool *pgxpool.Pool, chi chi.Router) {
	repo := postgres.New(pool)

	service, _ := identificationbus.New(repo, logger)

	identificationapp.Register(identificationapp.Config{
		Log:     logger,
		Service: service,
		Router:  chi,
		Repo:    repo,
	})
}
