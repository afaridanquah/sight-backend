package businessapp

import (
	"encoding/json"
	"net/http"

	"bitbucket.org/msafaridanquah/verifylab-service/business/domain/businessbus"
	"bitbucket.org/msafaridanquah/verifylab-service/foundation/ierr"
	"bitbucket.org/msafaridanquah/verifylab-service/foundation/logger"
	"bitbucket.org/msafaridanquah/verifylab-service/foundation/web"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
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
	ctx := r.Context()
	if err := json.NewDecoder(r.Body).Decode(&napp); err != nil {
		web.RenderErrorResponse(ctx, w, r, "invalid request",
			ierr.WrapErrorf(err, ierr.InvalidArgument, "json decoder"))
		return
	}

	defer func() {
		err := r.Body.Close()
		if err != nil {
			return
		}
	}()

	if err := napp.Validate(); err != nil {
		web.RenderErrorResponse(ctx, w, r, "validation failed", ierr.WrapErrorf(err, ierr.InvalidArgument, "json decoder"))
		return
	}

	newbus, err := toBusBusiness(napp)
	if err != nil {
		web.RenderErrorResponse(ctx, w, r, err.Error(), err)
		return
	}

	bbus, err := a.srv.Create(r.Context(), newbus)
	if err != nil {
		web.RenderErrorResponse(ctx, w, r, err.Error(), err)
		return
	}

	appBus := toAppBusiness(bbus)
	web.RenderResponse(http.StatusCreated, w, r, appBus)
}

func (a *App) update(w http.ResponseWriter, r *http.Request) {
	var up UpdateBusiness
	ctx := r.Context()
	id := chi.URLParam(r, "id")
	parsedBusinessID, err := uuid.Parse(id)
	if err != nil {
		web.RenderErrorResponse(ctx, w, r, err.Error(), err)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&up); err != nil {
		web.RenderErrorResponse(ctx, w, r, "invalid request",
			ierr.WrapErrorf(err, ierr.InvalidArgument, "json decoder"))
		return
	}

	defer func() {
		err := r.Body.Close()
		if err != nil {
			return
		}
	}()

	if err := up.Validate(); err != nil {
		web.RenderErrorResponse(ctx, w, r, "validation failed", ierr.WrapErrorf(err, ierr.InvalidArgument, "json decoder"))
		return
	}

	upbus, err := toBusUpdateBusiness(up)
	if err != nil {
		web.RenderErrorResponse(ctx, w, r, err.Error(), err)
		return
	}

	bus, err := a.srv.FindByID(ctx, parsedBusinessID)
	if err != nil {
		web.RenderErrorResponse(ctx, w, r, err.Error(), err)
		return
	}

	bbus, err := a.srv.Update(r.Context(), bus, upbus)
	if err != nil {
		web.RenderErrorResponse(ctx, w, r, err.Error(), err)
		return
	}

	appBus := toAppBusiness(bbus)
	web.RenderResponse(http.StatusCreated, w, r, appBus)
}

func (a *App) upload(w http.ResponseWriter, r *http.Request) {
	var napp NewDocument
	var ctx = r.Context()
	if err := json.NewDecoder(r.Body).Decode(&napp); err != nil {
		web.RenderErrorResponse(ctx, w, r, "invalid request",
			ierr.WrapErrorf(err, ierr.InvalidArgument, "json decoder"))
		return
	}

	defer func() {
		err := r.Body.Close()
		if err != nil {
			return
		}
	}()
}

func (a *App) findByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var ctx = r.Context()

	parsedID, err := uuid.Parse(id)
	if err != nil {
		web.RenderErrorResponse(ctx, w, r, err.Error(), err)
		return
	}

	bcus, err := a.srv.FindByID(r.Context(), parsedID)
	if err != nil {
		web.RenderErrorResponse(ctx, w, r, err.Error(), ierr.WrapErrorf(err, ierr.NotFound, "querybyidandbusinessid"))
		return
	}

	appBusiness := toAppBusiness(bcus)
	web.RenderResponse(http.StatusOK, w, r, appBusiness)
}
