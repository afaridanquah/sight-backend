package customerapp

import (
	"encoding/json"
	"net/http"

	"bitbucket.org/msafaridanquah/verifylab-service/business/domain/customerbus"
	"bitbucket.org/msafaridanquah/verifylab-service/foundation/ierr"
	"bitbucket.org/msafaridanquah/verifylab-service/foundation/logger"
	"bitbucket.org/msafaridanquah/verifylab-service/foundation/web"
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

	newbus, err := toBusinessNewCustomer(napp)
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

	appBalance := toAppCustomer(bcus)
	web.RenderResponse(w, r, appBalance, http.StatusCreated)
}

// func (a *App) customer(ctx context.Context, id string) (Customer, error) {
// 	bcus, err := a.srv.FindByID(ctx, id)

// 	if err != nil {
// 		return Customer{}, ierr.WrapErrorf(err, ierr.ErrorCodeInvalidArgument, "customerapp.QueryByID")
// 	}

// 	return toAppCustomer(bcus), nil
// }

// func (a *App) query(ctx context.Context) ([]Customer, error) {
// 	results, err := a.srv.All(ctx)

// 	if err != nil {
// 		return []Customer{}, ierr.WrapErrorf(err, ierr.ErrorCodeInvalidArgument, "customerapp.All")
// 	}

// 	return toAppCustomers(results), nil
// }
