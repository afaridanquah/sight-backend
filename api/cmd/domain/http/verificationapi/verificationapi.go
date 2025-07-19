package verificationapi

import (
	"bitbucket.org/msafaridanquah/verifylab-service/app/domain/verificationapp"
	"bitbucket.org/msafaridanquah/verifylab-service/business/domain/customerbus"
	"bitbucket.org/msafaridanquah/verifylab-service/business/domain/customerbus/postgres"
	"bitbucket.org/msafaridanquah/verifylab-service/business/domain/verificationbus"
	verificationPostgres "bitbucket.org/msafaridanquah/verifylab-service/business/domain/verificationbus/postgres"
	"bitbucket.org/msafaridanquah/verifylab-service/foundation/envvar"
	"bitbucket.org/msafaridanquah/verifylab-service/foundation/logger"
	"bitbucket.org/msafaridanquah/verifylab-service/foundation/vaulti"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Routes(logger *logger.Logger, pool *pgxpool.Pool, envvar *envvar.Configuration, vaulti *vaulti.Vaulty, chi chi.Router) {
	customerRepo := postgres.New(pool, pool, vaulti)
	verificationRepo := verificationPostgres.New(pool, pool, vaulti)

	customerService, _ := customerbus.New(customerRepo, logger)
	verificationService, _ := verificationbus.New(verificationRepo, logger)

	verificationapp.Register(verificationapp.Config{
		Log:                 logger,
		VerificationService: verificationService,
		CustomerService:     customerService,
		Router:              chi,
	})
}
