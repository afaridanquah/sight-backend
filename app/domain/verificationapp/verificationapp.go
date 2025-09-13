package verificationapp

import (
	"encoding/json"
	"net/http"

	"bitbucket.org/msafaridanquah/verifylab-service/business/domain/businessbus"
	"bitbucket.org/msafaridanquah/verifylab-service/business/domain/customerbus"
	"bitbucket.org/msafaridanquah/verifylab-service/business/domain/verificationbus"
	verbus_vo "bitbucket.org/msafaridanquah/verifylab-service/business/domain/verificationbus/valueobject"
	"bitbucket.org/msafaridanquah/verifylab-service/foundation/ierr"
	"bitbucket.org/msafaridanquah/verifylab-service/foundation/logger"
	"bitbucket.org/msafaridanquah/verifylab-service/foundation/web"
	"github.com/google/uuid"
)

type App struct {
	customerService    *customerbus.Service
	businessService    *businessbus.Service
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
	var ctx = r.Context()
	if err := json.NewDecoder(r.Body).Decode(&napp); err != nil {
		web.RenderErrorResponse(ctx, w, r, "invalid request",
			ierr.WrapErrorf(err, ierr.InvalidArgument, "json decoder"))
		return
	}

	defer func() {
		_ = r.Body.Close()
	}()

	if err := napp.Validate(); err != nil {
		app.log.Error(ctx, "validation", err)
		web.RenderErrorResponse(ctx, w, r, "validation", ierr.WrapErrorf(err, ierr.InvalidArgument, "json decoder"))
		return
	}

	switch {
	case napp.CustomerID != "":
		parsedCustomerID, err := uuid.Parse(napp.CustomerID)
		if err != nil {
			app.log.Error(ctx, "parse", err)
			web.RenderErrorResponse(ctx, w, r, "parse uuid", err)
			return
		}

		customer, err := app.customerService.FindByIDAndOrgID(r.Context(), parsedCustomerID)
		if err != nil {
			app.log.Error(ctx, "customerservice.findbyidandorgid", err)
			web.RenderErrorResponse(ctx, w, r, "customer not found", ierr.WrapErrorf(err, ierr.NotFound, ""))
			return
		}

		c, err := toBusVoCustomer(customer)
		if err != nil {
			app.log.Error(ctx, "tobusvocustomer", err)
			web.RenderErrorResponse(ctx, w, r, "tobusvocustomer", err)
			return
		}

		newbus, err := toBusNewCustomerVerification(napp, c)
		if err != nil {
			web.RenderErrorResponse(ctx, w, r, "customer not found", err)
		}

		vbus, err := app.verficationService.Create(r.Context(), newbus)
		if err != nil {
			web.RenderErrorResponse(ctx, w, r, "bus service create", err)
		}
		vapp := toAppVerification(vbus)
		web.RenderResponse(http.StatusOK, w, r, vapp)
		return

	case napp.BusinessID != "":
		parsedBusinessID, err := uuid.Parse(napp.BusinessID)
		if err != nil {
			app.log.Error(ctx, "parse", err)
			web.RenderErrorResponse(ctx, w, r, "parse uuid", err)
			return
		}

		business, err := app.businessService.FindByID(r.Context(), parsedBusinessID)
		if err != nil {
			app.log.Error(ctx, "businessService.findbyid", err)
			web.RenderErrorResponse(ctx, w, r, "business not found", ierr.WrapErrorf(err, ierr.NotFound, ""))
			return
		}

		countryCode := business.CountryOfIncorporation.Alpha2()
		b, err := verbus_vo.NewBusiness(business.ID, business.LegalName, &countryCode, &business.RegistrationNumber)
		if err != nil {
			web.RenderErrorResponse(ctx, w, r, "newbusiness", err)
			return
		}

		newbus, err := toBusNewBusinessVerification(napp, b)
		if err != nil {
			web.RenderErrorResponse(ctx, w, r, "tobusnewbusinessverification", err)
		}

		vbus, err := app.verficationService.Create(r.Context(), newbus)
		if err != nil {
			web.RenderErrorResponse(ctx, w, r, "verficationservice.create", err)
		}
		vapp := toAppVerification(vbus)
		web.RenderResponse(http.StatusOK, w, r, vapp)
		return

	default:
		return
	}
}
