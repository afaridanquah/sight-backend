package countryapp

import (
	"net/http"

	"bitbucket.org/msafaridanquah/sight-backend/foundation/logger"
)

type App struct {
	log *logger.Logger
}

func newApp(log *logger.Logger) *App {
	return &App{
		log: log,
	}
}

func (app *App) index(w http.ResponseWriter, r *http.Request) {
}
