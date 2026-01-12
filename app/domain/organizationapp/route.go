package organizationapp

import (
	"bitbucket.org/msafaridanquah/sight-backend/business/domain/organizationbus"
	"bitbucket.org/msafaridanquah/sight-backend/business/domain/organizationbus/postgres"
	"bitbucket.org/msafaridanquah/sight-backend/foundation/envvar"
	"bitbucket.org/msafaridanquah/sight-backend/foundation/logger"
	"github.com/go-chi/chi/v5"
)

type Config struct {
	Log        *logger.Logger
	Router     chi.Router
	OrgService *organizationbus.Service
	Repo       *postgres.Repository
	EnvVar     *envvar.Configuration
}

func Register(conf Config) {
	app := newApp(conf.OrgService, conf.Log, conf.EnvVar)

	conf.Router.Post("/organizations", app.create)
}
