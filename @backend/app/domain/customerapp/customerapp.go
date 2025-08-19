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
	var ctx = r.Context()
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
		a.log.Error(ctx, "customerapp.validate", err)
		web.RenderErrorResponse(a.log, w, r, "validation failed", ierr.WrapErrorf(err, ierr.ErrorCodeInvalidArgument, "json decoder"))
		return
	}

	newbus, err := toBusNewCustomer(napp)
	a.log.Info(ctx, "customerapp.toBusNewCustomer", newbus)
	if err != nil {
		web.RenderErrorResponse(a.log, w, r, err.Error(), err)
		return
	}

	bcus, err := a.srv.Create(r.Context(), newbus)
	if err != nil {
		web.RenderErrorResponse(a.log, w, r, err.Error(), err)
		return
	}

	appCustomer := toAppCustomer(bcus)
	web.RenderResponse(http.StatusCreated, w, r, appCustomer)
}

func (a *App) customer(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	parsedCustomerID, err := uuid.Parse(id)
	if err != nil {
		web.RenderErrorResponse(a.log, w, r, err.Error(), err)
		return
	}

	bcus, err := a.srv.QueryByIDAndBusinessID(r.Context(), parsedCustomerID)
	if err != nil {
		web.RenderErrorResponse(a.log, w, r, err.Error(), err)
		return
	}

	appCustomer := toAppCustomer(bcus)
	web.RenderResponse(http.StatusCreated, w, r, appCustomer)
}
