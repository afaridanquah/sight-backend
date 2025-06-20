package customerapi

import (
	"bitbucket.org/msafaridanquah/verifylab-service/app/domain/customerapp"
	"bitbucket.org/msafaridanquah/verifylab-service/business/domain/customerbus"
	"bitbucket.org/msafaridanquah/verifylab-service/business/domain/customerbus/postgres"
	"bitbucket.org/msafaridanquah/verifylab-service/foundation/logger"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Routes(logger *logger.Logger, pool *pgxpool.Pool, chi chi.Router) {
	repo := postgres.New(pool, pool)

	service, _ := customerbus.New(repo, logger)

	customerapp.Register(customerapp.Config{
		Log:     logger,
		Service: service,
		Router:  chi,
		Repo:    repo,
	})
}
