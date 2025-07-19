package verificationapp

import (
	"bitbucket.org/msafaridanquah/verifylab-service/business/domain/customerbus"
	"bitbucket.org/msafaridanquah/verifylab-service/business/domain/verificationbus"
	"bitbucket.org/msafaridanquah/verifylab-service/business/domain/verificationbus/postgres"
	"bitbucket.org/msafaridanquah/verifylab-service/foundation/logger"
	"github.com/go-chi/chi/v5"
)

type Config struct {
	Log                 *logger.Logger
	Router              chi.Router
	VerificationService *verificationbus.Service
	CustomerService     *customerbus.Service
	Repo                *postgres.Repository
}

func Register(conf Config) {
	app := newApp(conf.VerificationService, conf.CustomerService, conf.Log)

	conf.Router.Post("/screening", app.screen)
}
