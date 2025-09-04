package businessapp

import (
	"bitbucket.org/msafaridanquah/verifylab-service/business/domain/businessbus"
	"bitbucket.org/msafaridanquah/verifylab-service/business/domain/businessbus/postgres"
	"bitbucket.org/msafaridanquah/verifylab-service/foundation/logger"
	"github.com/go-chi/chi/v5"
)

type Config struct {
	Log     *logger.Logger
	Router  chi.Router
	Service *businessbus.Service
	Repo    *postgres.Repository
}

func Register(conf Config) {
	app := newApp(conf.Service, conf.Log)

	conf.Router.Get("/businesses/{id}", app.findByID)
	conf.Router.Put("/businesses/{id}", app.update)
	conf.Router.Post("/businesses", app.create)
	conf.Router.Post("/businesses/{id}/documents", app.upload)

}
