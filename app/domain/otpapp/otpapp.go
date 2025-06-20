package otpapp

import (
	"encoding/json"
	"net/http"

	"bitbucket.org/msafaridanquah/verifylab-service/business/domain/otpbus"
	"bitbucket.org/msafaridanquah/verifylab-service/foundation/ierr"
	"bitbucket.org/msafaridanquah/verifylab-service/foundation/logger"
	"bitbucket.org/msafaridanquah/verifylab-service/foundation/web"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type App struct {
	log *logger.Logger
	srv *otpbus.Service
}

func newApp(srv *otpbus.Service, log *logger.Logger) *App {
	return &App{
		srv: srv,
		log: log,
	}
}

func (app *App) create(w http.ResponseWriter, r *http.Request) {
	var napp NewOTP

	if err := json.NewDecoder(r.Body).Decode(&napp); err != nil {
		web.RenderErrorResponse(w, r, "invalid request",
			ierr.WrapErrorf(err, ierr.ErrorCodeInvalidArgument, "json decoder"))
		return
	}

	id := chi.URLParam(r, "id")
	custID, err := uuid.Parse(id)
	if err != nil {
		app.log.Error(r.Context(), "otpapp.parseCustID", err)
		web.RenderErrorResponse(w, r, err.Error(), err)
		return
	}

	defer r.Body.Close()

	if err := napp.Validate(); err != nil {
		app.log.Error(r.Context(), "otpapp.validate", err)
		web.RenderErrorResponse(w, r, "validation failed", ierr.WrapErrorf(err, ierr.ErrorCodeInvalidArgument, "json decoder"))
		return
	}

	newbus, err := toBusNewOTP(napp)
	if err != nil {
		app.log.Error(r.Context(), "otpapp.toBusNewOTP", err)
		web.RenderErrorResponse(w, r, err.Error(), err)
		return
	}

	bcus, err := app.srv.Create(r.Context(), custID, newbus)
	if err != nil {
		app.log.Error(r.Context(), "otpapp.srv.Create", err)

		web.RenderErrorResponse(w, r, err.Error(), err)
		return
	}

	appBalance := toAppOTP(bcus)
	web.RenderResponse(w, r, appBalance, http.StatusCreated)
}
