package customerapp

import (
	"encoding/json"
	"net/http"

	"bitbucket.org/msafaridanquah/verifylab-service/business/domain/customerbus"
	"bitbucket.org/msafaridanquah/verifylab-service/foundation/ierr"
	"bitbucket.org/msafaridanquah/verifylab-service/foundation/logger"
	"bitbucket.org/msafaridanquah/verifylab-service/foundation/web"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type App struct {
	srv *customerbus.Service
	log *logger.Logger
}

func newApp(srv *customerbus.Service, log *logger.Logger) *App {
	return &App{
		srv: srv,
		log: log,
	}
}

func (a *App) create(w http.ResponseWriter, r *http.Request) {
	var napp NewCustomer
	if err := json.NewDecoder(r.Body).Decode(&napp); err != nil {
		web.RenderErrorResponse(w, r, "invalid request",
			ierr.WrapErrorf(err, ierr.ErrorCodeInvalidArgument, "json decoder"))
		return
	}

	defer r.Body.Close()

	if err := napp.Validate(); err != nil {
		a.log.Error(r.Context(), "customerapp.validate", err)
		web.RenderErrorResponse(w, r, "validation failed", ierr.WrapErrorf(err, ierr.ErrorCodeInvalidArgument, "json decoder"))
		return
	}

	newbus, err := toBusNewCustomer(napp)
	a.log.Info(r.Context(), "customerapp.toBusNewCustomer", newbus)
	if err != nil {
		a.log.Error(r.Context(), "customerapp.toBusNewCustomer", err)
		web.RenderErrorResponse(w, r, err.Error(), err)
		return
	}

	bcus, err := a.srv.Create(r.Context(), newbus)
	if err != nil {
		a.log.Error(r.Context(), "customerapp.srv.Create", err)

		web.RenderErrorResponse(w, r, err.Error(), err)
		return
	}

	appCustomer := toAppCustomer(bcus)
	web.RenderResponse(w, r, appCustomer, http.StatusCreated)
}

func (a *App) customer(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	parsedCustomerID, err := uuid.Parse(id)

	bcus, err := a.srv.QueryByIDAndBusinessID(r.Context(), parsedCustomerID)
	if err != nil {
		web.RenderErrorResponse(w, r, err.Error(), err)
		return
	}

	appCustomer := toAppCustomer(bcus)
	web.RenderResponse(w, r, appCustomer, http.StatusCreated)
}
