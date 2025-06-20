package otpapi

import (
	"bitbucket.org/msafaridanquah/verifylab-service/app/domain/otpapp"
	"bitbucket.org/msafaridanquah/verifylab-service/business/domain/otpbus"
	"bitbucket.org/msafaridanquah/verifylab-service/business/domain/otpbus/postgres"
	"bitbucket.org/msafaridanquah/verifylab-service/foundation/logger"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Routes(logger *logger.Logger, pool *pgxpool.Pool, chi chi.Router) {
	repo := postgres.New(pool, pool)

	service, _ := otpbus.New(repo, logger)

	otpapp.Register(otpapp.Config{
		Log:     logger,
		Service: service,
		Router:  chi,
		Repo:    repo,
	})
}
