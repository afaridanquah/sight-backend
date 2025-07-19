package verificationapp

import (
	"encoding/json"
	"net/http"

	"bitbucket.org/msafaridanquah/verifylab-service/business/domain/customerbus"
	"bitbucket.org/msafaridanquah/verifylab-service/business/domain/verificationbus"
	"bitbucket.org/msafaridanquah/verifylab-service/foundation/ierr"
	"bitbucket.org/msafaridanquah/verifylab-service/foundation/logger"
	"bitbucket.org/msafaridanquah/verifylab-service/foundation/web"
	"github.com/google/uuid"
)

type App struct {
	customerService    *customerbus.Service
	verficationService *verificationbus.Service
	log                *logger.Logger
}

func newApp(vs *verificationbus.Service, cs *customerbus.Service, log *logger.Logger) *App {
	return &App{
		customerService:    cs,
		verficationService: vs,
		log:                log,
	}
}

func (app *App) screen(w http.ResponseWriter, r *http.Request) {
	var napp NewVerification
	if err := json.NewDecoder(r.Body).Decode(&napp); err != nil {
		web.RenderErrorResponse(w, r, "invalid request",
			ierr.WrapErrorf(err, ierr.ErrorCodeInvalidArgument, "json decoder"))
		return
	}

	defer r.Body.Close()

	if err := napp.Validate(); err != nil {
		app.log.Error(r.Context(), "verificationapp.validate", err)
		web.RenderErrorResponse(w, r, "validation", ierr.WrapErrorf(err, ierr.ErrorCodeInvalidArgument, "json decoder"))
		return
	}

	parsedCustomerID, err := uuid.Parse(napp.CustomerID)
	if err != nil {
		web.RenderErrorResponse(w, r, "parse uuid", err)
	}

	customer, err := app.customerService.QueryByIDAndBusinessID(r.Context(), parsedCustomerID)
	if err != nil {
		web.RenderErrorResponse(w, r, "customer not found", err)
	}

	voCustomer, err := toBusVoCustomer(customer)
	if err != nil {
		web.RenderErrorResponse(w, r, "tobusvocustomer", err)
	}

	newbus, err := toBusNewVerification(napp, voCustomer)
	if err != nil {
		app.log.Error(r.Context(), "tobusnewverification", err)
		web.RenderErrorResponse(w, r, "customer not found", err)
	}

	vbus, err := app.verficationService.Create(r.Context(), newbus)
	if err != nil {
		app.log.Error(r.Context(), "Create", err)
		web.RenderErrorResponse(w, r, "bus service create", err)
	}

	vapp := toAppVerification(vbus)

	web.RenderResponse(w, r, vapp, http.StatusOK)
}
