package otpapi

import (
	"bitbucket.org/msafaridanquah/sight-backend/app/domain/otpapp"
	"bitbucket.org/msafaridanquah/sight-backend/business/domain/customerbus"
	cp "bitbucket.org/msafaridanquah/sight-backend/business/domain/customerbus/postgres"
	"bitbucket.org/msafaridanquah/sight-backend/business/domain/otpbus"
	op "bitbucket.org/msafaridanquah/sight-backend/business/domain/otpbus/postgres"
	"bitbucket.org/msafaridanquah/sight-backend/foundation/envvar"
	"bitbucket.org/msafaridanquah/sight-backend/foundation/logger"
	"bitbucket.org/msafaridanquah/sight-backend/foundation/vaulti"
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
