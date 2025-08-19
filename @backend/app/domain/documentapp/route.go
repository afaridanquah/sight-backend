package documentapp

import (
	"bitbucket.org/msafaridanquah/verifylab-service/business/domain/documentbus"
	// "bitbucket.org/msafaridanquah/verifylab-service/business/domain/documentbus/postgres"
	"bitbucket.org/msafaridanquah/verifylab-service/foundation/logger"
	"github.com/go-chi/chi/v5"
)

type Config struct {
	Log     *logger.Logger
	Router  chi.Router
	Service *documentbus.Service
	// Repo    *postgres.Repository
}

func Register(conf Config) {
	app := newApp(conf.Service, conf.Log)

	// conf.Router.Post("/customers", app.create)
	// conf.Router.Get("/customers/{id}", app.customer)
}
