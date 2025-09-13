package customerapp

import (
	"bitbucket.org/msafaridanquah/verifylab-service/business/domain/customerbus"
	"bitbucket.org/msafaridanquah/verifylab-service/business/domain/customerbus/postgres"
	"bitbucket.org/msafaridanquah/verifylab-service/business/domain/documentbus"
	"bitbucket.org/msafaridanquah/verifylab-service/foundation/envvar"
	"bitbucket.org/msafaridanquah/verifylab-service/foundation/logger"
	"github.com/go-chi/chi/v5"
)

type Config struct {
	Log             *logger.Logger
	Router          chi.Router
	Service         *customerbus.Service
	DocumentService *documentbus.Service
	Repo            *postgres.Repository
	EnvVar          *envvar.Configuration
}

func Register(conf Config) {
	app := newApp(conf.Service, conf.DocumentService, conf.Log, conf.EnvVar)

	conf.Router.Post("/customers/{id}/upload", app.upload)
	conf.Router.Post("/customers", app.create)
	conf.Router.Get("/customers/{id}", app.findByID)
}
