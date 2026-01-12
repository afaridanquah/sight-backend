package organizationapp

import (
	"encoding/json"
	"net/http"

	"bitbucket.org/msafaridanquah/sight-backend/app/sdk/errs"
	"bitbucket.org/msafaridanquah/sight-backend/business/domain/organizationbus"
	"bitbucket.org/msafaridanquah/sight-backend/foundation/envvar"
	"bitbucket.org/msafaridanquah/sight-backend/foundation/logger"
	"bitbucket.org/msafaridanquah/sight-backend/foundation/web"
)

type App struct {
	orgService *organizationbus.Service
	log        *logger.Logger
	envvar     *envvar.Configuration
}

func newApp(orgService *organizationbus.Service, log *logger.Logger, envvar *envvar.Configuration) *App {
	return &App{
		orgService: orgService,
		log:        log,
		envvar:     envvar,
	}
}

func (app *App) create(w http.ResponseWriter, r *http.Request) {
	var napp NewOrganization

	var ctx = r.Context()
	if err := json.NewDecoder(r.Body).Decode(&napp); err != nil {
		web.RenderErrorResponse(ctx, w, r, "invalid request",
			errs.New(errs.Internal, err))
		return
	}

	defer func() {
		_ = r.Body.Close()
	}()

	if err := napp.Validate(); err != nil {
		web.RenderErrorResponse(ctx, w, r,
			"validation failed", errs.New(errs.InvalidArgument, err))
		return
	}

	nbus, err := toBusNewOrganization(napp)
	if err != nil {
		web.RenderErrorResponse(ctx, w, r,
			"validation failed", errs.New(errs.InvalidArgument, err))
		return
	}

	bus, err := app.orgService.Create(ctx, nbus)
	if err != nil {
		web.RenderErrorResponse(ctx, w, r,
			"validation failed", errs.Newf(errs.Internal, "create: bus[%+v]: %s", bus, err))
		return
	}

	org := toAppOrganization(bus)
	web.RenderResponse(http.StatusCreated, w, r, org)
}
