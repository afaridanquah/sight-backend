package govapp

import (
	"bitbucket.org/msafaridanquah/sight-backend/business/domain/govbus"
	"bitbucket.org/msafaridanquah/sight-backend/foundation/envvar"
	"bitbucket.org/msafaridanquah/sight-backend/foundation/logger"
)

type App struct {
	srv    *govbus.Service
	log    *logger.Logger
	envvar *envvar.Configuration
}

func newApp(srv *govbus.Service, logger *logger.Logger, envvar *envvar.Configuration) *App {
	return &App{
		srv:    srv,
		log:    logger,
		envvar: envvar,
	}
}
