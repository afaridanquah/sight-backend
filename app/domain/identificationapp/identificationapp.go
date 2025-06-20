package identificationapp

import (
	"encoding/json"
	"net/http"

	"bitbucket.org/msafaridanquah/verifylab-service/business/domain/identificationbus"
	"bitbucket.org/msafaridanquah/verifylab-service/foundation/ierr"
	"bitbucket.org/msafaridanquah/verifylab-service/foundation/logger"
	"bitbucket.org/msafaridanquah/verifylab-service/foundation/web"
)

type App struct {
	srv *identificationbus.Service
	log *logger.Logger
}

func newApp(srv *identificationbus.Service, log *logger.Logger) *App {
	return &App{
		srv: srv,
		log: log,
	}
}

func (app *App) create(w http.ResponseWriter, r *http.Request) {
	var napp NewIdentification
	if err := json.NewDecoder(r.Body).Decode(&napp); err != nil {
		web.RenderErrorResponse(w, r, "invalid request",
			ierr.WrapErrorf(err, ierr.ErrorCodeInvalidArgument, "json decoder"))
		return
	}

	defer r.Body.Close()

	if err := napp.Validate(); err != nil {
		app.log.Error(r.Context(), "identificationapp.validate", err)
		web.RenderErrorResponse(w, r, "validation failed", ierr.WrapErrorf(err, ierr.ErrorCodeInvalidArgument, "json decoder"))
		return
	}

	newbus, err := toBusNewIdentification(napp)
	if err != nil {
		app.log.Error(r.Context(), "customerapp.toBusNewCustomer", err)
		web.RenderErrorResponse(w, r, err.Error(), err)
		return
	}

	bid, err := app.srv.Create(r.Context(), newbus)
	if err != nil {
		app.log.Error(r.Context(), "customerapp.srv.Create", err)

		web.RenderErrorResponse(w, r, err.Error(), err)
		return
	}

	appBalance := toAppIdentification(bid)
	web.RenderResponse(w, r, appBalance, http.StatusCreated)
}
