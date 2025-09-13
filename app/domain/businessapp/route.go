package businessapp

import (
	"bitbucket.org/msafaridanquah/sight-backend/business/domain/businessbus"
	"bitbucket.org/msafaridanquah/sight-backend/business/domain/businessbus/postgres"
	"bitbucket.org/msafaridanquah/sight-backend/foundation/logger"
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

	conf.Router.Delete("/businesses/{id}", app.delete)
	conf.Router.Get("/businesses/{id}", app.findByID)
	conf.Router.Put("/businesses/{id}", app.update)
	conf.Router.Post("/businesses", app.create)
	conf.Router.Post("/businesses/{id}/documents", app.upload)
}
