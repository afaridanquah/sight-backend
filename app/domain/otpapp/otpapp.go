package otpapp

import (
	"encoding/json"
	"fmt"
	"net/http"

	"bitbucket.org/msafaridanquah/verifylab-service/business/domain/customerbus"
	"bitbucket.org/msafaridanquah/verifylab-service/business/domain/otpbus"
	"bitbucket.org/msafaridanquah/verifylab-service/foundation/ierr"
	"bitbucket.org/msafaridanquah/verifylab-service/foundation/logger"
	"bitbucket.org/msafaridanquah/verifylab-service/foundation/web"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type App struct {
	log             *logger.Logger
	customerService *customerbus.Service
	otpService      *otpbus.Service
}

func newApp(os *otpbus.Service, cs *customerbus.Service, log *logger.Logger) *App {
	return &App{
		otpService:      os,
		customerService: cs,
		log:             log,
	}
}

func (app *App) create(w http.ResponseWriter, r *http.Request) {
	var napp NewOTP
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

	id := chi.URLParam(r, "id")
	custID, err := uuid.Parse(id)
	if err != nil {
		app.log.Error(r.Context(), "otpapp.parseCustID", err)
		web.RenderErrorResponse(ctx, w, r, err.Error(), err)
		return
	}

	customer, err := app.customerService.FindByIDAndOrgID(r.Context(), custID)
	if err != nil {
		web.RenderErrorResponse(ctx, w, r, err.Error(), err)
		return
	}

	switch {
	case napp.Channel == "PHONE" && customer.PhoneNumber.IsEmpty():
		web.RenderErrorResponse(ctx, w, r, "validation", ierr.WrapErrorf(fmt.Errorf("customer: %s has no valid phone number", customer.ID), ierr.InvalidArgument, "empty phone number"))
		return
	case napp.Channel == "EMAIL" && customer.Email.IsEmpty():
		web.RenderErrorResponse(ctx, w, r, "validation", ierr.WrapErrorf(fmt.Errorf("customer: %s has no valid email address", customer.ID), ierr.InvalidArgument, "empty phone number"))
		return
	}

	if err := napp.Validate(); err != nil {
		app.log.Error(r.Context(), "otpapp.validate", err)
		web.RenderErrorResponse(ctx, w, r, "validation", ierr.WrapErrorf(err, ierr.InvalidArgument, "empty email"))
		return
	}

	newbus, err := toBusNewOTP(napp)
	if err != nil {
		app.log.Error(r.Context(), "otpapp.toBusNewOTP", err)
		web.RenderErrorResponse(ctx, w, r, err.Error(), err)
		return
	}

	if napp.Channel == "PHONE" {
		newbus.Destination = customer.PhoneNumber.E164Format
	}

	if napp.Channel == "EMAIL" {
		newbus.Destination = customer.Email.String()
	}

	bcus, err := app.otpService.Create(r.Context(), custID, newbus)
	if err != nil {
		app.log.Error(r.Context(), "otpapp.srv.Create", err)

		web.RenderErrorResponse(ctx, w, r, err.Error(), err)
		return
	}

	appBalance := toAppOTP(bcus)
	web.RenderResponse(http.StatusCreated, w, r, appBalance)
}

func (app *App) verify(w http.ResponseWriter, r *http.Request) {
	var napp VerifyOTP
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

	id := chi.URLParam(r, "id")
	custID, err := uuid.Parse(id)
	if err != nil {
		app.log.Error(ctx, "otpapp.parseCustID", err)
		web.RenderErrorResponse(ctx, w, r, err.Error(), err)
		return
	}

	customer, err := app.customerService.FindByIDAndOrgID(ctx, custID)
	if err != nil {
		web.RenderErrorResponse(ctx, w, r, err.Error(), err)
		return
	}

	if err := napp.Validate(); err != nil {
		app.log.Error(ctx, "otpapp.validate", err)
		web.RenderErrorResponse(ctx, w, r, "validation", ierr.WrapErrorf(err, ierr.InvalidArgument, "empty email"))
		return
	}

	bus, err := toBusVerifyOTP(napp, customer.ID)
	if err != nil {
		web.RenderErrorResponse(ctx, w, r, err.Error(), err)
		return
	}

	otpbus, err := app.otpService.Verify(ctx, bus)
	if err != nil {
		web.RenderErrorResponse(ctx, w, r, err.Error(), err)
		return
	}

	appBalance := toAppOTP(otpbus)
	web.RenderResponse(http.StatusOK, w, r, appBalance)
}
