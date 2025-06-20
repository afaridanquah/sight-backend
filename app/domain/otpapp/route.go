package otpapp

import (
	"bitbucket.org/msafaridanquah/verifylab-service/business/domain/otpbus"
	"bitbucket.org/msafaridanquah/verifylab-service/business/domain/otpbus/postgres"
	"bitbucket.org/msafaridanquah/verifylab-service/foundation/logger"
	"github.com/go-chi/chi/v5"
)

type Config struct {
	Log     *logger.Logger
	Router  chi.Router
	Service *otpbus.Service
	Repo    *postgres.Repository
}

func Register(conf Config) {
	app := newApp(conf.Service, conf.Log)

	conf.Router.Post("/customers/{id}/otps", app.create)
	conf.Router.Put("/customers/{id}/otps", app.create)
}
