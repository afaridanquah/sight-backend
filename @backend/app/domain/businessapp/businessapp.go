package businessapp

import (
	"encoding/json"
	"net/http"

	"bitbucket.org/msafaridanquah/verifylab-service/business/domain/businessbus"
	"bitbucket.org/msafaridanquah/verifylab-service/foundation/ierr"
	"bitbucket.org/msafaridanquah/verifylab-service/foundation/logger"
	"bitbucket.org/msafaridanquah/verifylab-service/foundation/web"
)

type App struct {
	srv *businessbus.Service
	log *logger.Logger
}

func newApp(srv *businessbus.Service, log *logger.Logger) *App {
	return &App{
		srv: srv,
		log: log,
	}
}

func (a *App) create(w http.ResponseWriter, r *http.Request) {
	var napp NewBusiness
	if err := json.NewDecoder(r.Body).Decode(&napp); err != nil {
		web.RenderErrorResponse(a.log, w, r, "invalid request",
			ierr.WrapErrorf(err, ierr.ErrorCodeInvalidArgument, "json decoder"))
		return
	}

	defer func() {
		err := r.Body.Close()
		if err != nil {
			return
		}
	}()

	if err := napp.Validate(); err != nil {
		web.RenderErrorResponse(a.log, w, r, "validation failed", ierr.WrapErrorf(err, ierr.ErrorCodeInvalidArgument, "json decoder"))
		return
	}

	newbus, err := toBusBusiness(napp)
	if err != nil {
		web.RenderErrorResponse(a.log, w, r, err.Error(), err)
		return
	}

	bbus, err := a.srv.Create(r.Context(), newbus)
	if err != nil {
		web.RenderErrorResponse(a.log, w, r, err.Error(), err)
		return
	}

	appBus := toAppBusiness(bbus)
	web.RenderResponse(http.StatusCreated, w, r, appBus)
}
