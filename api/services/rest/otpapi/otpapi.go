package otpapi

import (
	"bitbucket.org/msafaridanquah/verifylab-service/app/domain/otpapp"
	"bitbucket.org/msafaridanquah/verifylab-service/business/domain/customerbus"
	cp "bitbucket.org/msafaridanquah/verifylab-service/business/domain/customerbus/postgres"
	"bitbucket.org/msafaridanquah/verifylab-service/business/domain/otpbus"
	op "bitbucket.org/msafaridanquah/verifylab-service/business/domain/otpbus/postgres"
	"bitbucket.org/msafaridanquah/verifylab-service/foundation/envvar"
	"bitbucket.org/msafaridanquah/verifylab-service/foundation/logger"
	"bitbucket.org/msafaridanquah/verifylab-service/foundation/vaulti"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Routes(logger *logger.Logger, pool *pgxpool.Pool, envvar *envvar.Configuration, vaulti *vaulti.Vaulty, chi chi.Router) {
	cr := cp.New(pool, pool, vaulti)
	or := op.New(pool, pool, vaulti)

	cs, _ := customerbus.New(cr, logger)
	ts, _ := otpbus.New(or, logger)

	otpapp.Register(otpapp.Config{
		Log:             logger,
		Router:          chi,
		CustomerService: cs,
		OtpService:      ts,
	})
}
