package customerapp

import (
	"bitbucket.org/msafaridanquah/sight-backend/business/domain/customerbus"
	"bitbucket.org/msafaridanquah/sight-backend/business/domain/customerbus/postgres"
	"bitbucket.org/msafaridanquah/sight-backend/business/domain/documentbus"
	"bitbucket.org/msafaridanquah/sight-backend/business/domain/identificationbus"
	"bitbucket.org/msafaridanquah/sight-backend/foundation/envvar"
	"bitbucket.org/msafaridanquah/sight-backend/foundation/logger"
	"github.com/go-chi/chi/v5"
)

type Config struct {
	Log                   *logger.Logger
	Router                chi.Router
	Service               *customerbus.Service
	DocumentService       *documentbus.Service
	IdentificationService *identificationbus.Service
	Repo                  *postgres.Repository
	EnvVar                *envvar.Configuration
}

func Register(conf Config) {
	app := newApp(conf.Service, conf.DocumentService, conf.IdentificationService, conf.Log, conf.EnvVar)

	conf.Router.Get("/customers", app.query)
	conf.Router.Post("/customers/{id}/upload", app.upload)
	conf.Router.Post("/customers", app.create)
	conf.Router.Put("/customers/{id}", app.update)
	conf.Router.Get("/customers/{id}", app.findByID)
	conf.Router.Post("/customers/{id}/identifications", app.createIdentifications)
}
