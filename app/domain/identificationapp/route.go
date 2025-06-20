package identificationapp

import (
	"bitbucket.org/msafaridanquah/verifylab-service/business/domain/identificationbus"
	"bitbucket.org/msafaridanquah/verifylab-service/business/domain/identificationbus/postgres"
	"bitbucket.org/msafaridanquah/verifylab-service/foundation/logger"
	"github.com/go-chi/chi/v5"
)

type Config struct {
	Log     *logger.Logger
	Router  chi.Router
	Service *identificationbus.Service
	Repo    *postgres.Repository
}

func Register(conf Config) {
	app := newApp(conf.Service, conf.Log)

	conf.Router.Post("/identifications", app.create)
}
