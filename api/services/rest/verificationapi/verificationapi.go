package verificationapi

import (
	"bitbucket.org/msafaridanquah/sight-backend/app/domain/verificationapp"
	"bitbucket.org/msafaridanquah/sight-backend/business/domain/customerbus"
	"bitbucket.org/msafaridanquah/sight-backend/business/domain/customerbus/postgres"
	"bitbucket.org/msafaridanquah/sight-backend/business/domain/verificationbus"
	vp "bitbucket.org/msafaridanquah/sight-backend/business/domain/verificationbus/postgres"
	"bitbucket.org/msafaridanquah/sight-backend/foundation/envvar"
	"bitbucket.org/msafaridanquah/sight-backend/foundation/logger"
	"bitbucket.org/msafaridanquah/sight-backend/foundation/vaulti"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Routes(logger *logger.Logger, pool *pgxpool.Pool, envvar *envvar.Configuration, vaulti *vaulti.Vaulty, chi chi.Router) {
	customerRepo := postgres.New(pool, pool, vaulti)
	verificationRepo := vp.New(pool, pool, vaulti)

	customerService, _ := customerbus.New(customerRepo, logger)
	verificationService, _ := verificationbus.New(verificationRepo, logger)

	verificationapp.Register(verificationapp.Config{
		Log:                 logger,
		VerificationService: verificationService,
		CustomerService:     customerService,
		Router:              chi,
	})
}
