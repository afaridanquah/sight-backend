package customerapp

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"bitbucket.org/msafaridanquah/sight-backend/app/sdk/errs"
	"bitbucket.org/msafaridanquah/sight-backend/business/domain/customerbus"
	cvo "bitbucket.org/msafaridanquah/sight-backend/business/domain/customerbus/valueobject"
	"bitbucket.org/msafaridanquah/sight-backend/business/domain/documentbus"
	dvo "bitbucket.org/msafaridanquah/sight-backend/business/domain/documentbus/valueobject"
	"bitbucket.org/msafaridanquah/sight-backend/business/domain/identificationbus"
	"bitbucket.org/msafaridanquah/sight-backend/foundation/envvar"
	"bitbucket.org/msafaridanquah/sight-backend/foundation/logger"
	"bitbucket.org/msafaridanquah/sight-backend/foundation/web"

	"github.com/go-chi/chi/v5"
)

type App struct {
	srv                   *customerbus.Service
	documentService       *documentbus.Service
	identificationService *identificationbus.Service
	log                   *logger.Logger
	envvar                *envvar.Configuration
}

func newApp(srv *customerbus.Service, ds *documentbus.Service, is *identificationbus.Service, log *logger.Logger, envvar *envvar.Configuration) *App {
	return &App{
		srv:                   srv,
		log:                   log,
		documentService:       ds,
		envvar:                envvar,
		identificationService: is,
	}
}

func (a *App) query(w http.ResponseWriter, r *http.Request) {
	var ctx = r.Context()

	bcuss, err := a.srv.QueryByOrgID(r.Context())
	if err != nil {
		web.RenderErrorResponse(ctx, w, r, err.Error(), errs.New(errs.NotFound, err))
		return
	}

	customers := make([]Customer, len(bcuss))
	for k, v := range bcuss {
		appCustomer := toAppCustomer(v)
		customers[k] = appCustomer
	}

	web.RenderResponse(http.StatusOK, w, r, customers)
}

func (a *App) create(w http.ResponseWriter, r *http.Request) {
	var napp NewCustomer
	var ctx = r.Context()
	if err := json.NewDecoder(r.Body).Decode(&napp); err != nil {
		web.RenderErrorResponse(ctx, w, r, "invalid request",
			errs.New(errs.InvalidArgument, err))
		return
	}

	defer func() {
		err := r.Body.Close()
		if err != nil {
			return
		}
	}()

	if err := napp.Validate(); err != nil {
		web.RenderErrorResponse(ctx, w, r, "validation failed", errs.New(errs.InvalidArgument, err))
		return
	}

	newbus, err := toBusNewCustomer(napp)
	if err != nil {
		web.RenderErrorResponse(ctx, w, r, err.Error(), err)
		return
	}

	bcus, err := a.srv.Create(ctx, newbus)
	if err != nil {
		a.log.Info(ctx, "customerapp.srv.create", err)
		web.RenderErrorResponse(ctx, w, r, err.Error(), err)
		return
	}

	cus := toAppCustomer(bcus)

	web.RenderResponse(http.StatusCreated, w, r, cus)
}

func (a *App) update(w http.ResponseWriter, r *http.Request) {
	var ctx = r.Context()
	id := chi.URLParam(r, "id")

	parsedCustomerID, err := cvo.ParseCustomerID(id)
	if err != nil {
		web.RenderErrorResponse(ctx, w, r, err.Error(), err)
		return
	}

	var uapp UpdateCustomer
	if err := json.NewDecoder(r.Body).Decode(&uapp); err != nil {
		web.RenderErrorResponse(ctx, w, r, "invalid request",
			errs.New(errs.InvalidArgument, err))
		return
	}

	defer func() {
		err := r.Body.Close()
		if err != nil {
			return
		}
	}()

	if err := uapp.Validate(); err != nil {
		web.RenderErrorResponse(ctx, w, r, "validation failed", errs.New(errs.InvalidArgument, err))
		return
	}

	cus, err := a.srv.FindByIDAndOrgID(r.Context(), parsedCustomerID)
	if err != nil {
		web.RenderErrorResponse(ctx, w, r, err.Error(), errs.New(errs.NotFound, err))
		return
	}

	ubus, err := toBusUpdateCustomer(uapp)
	if err != nil {
		a.log.Error(ctx, "tobusupdatecustomer", err)
		web.RenderErrorResponse(ctx, w, r, err.Error(), errs.New(errs.Internal, err))
		return
	}

	bus, err := a.srv.Update(ctx, cus, ubus)
	if err != nil {
		web.RenderErrorResponse(ctx, w, r, err.Error(), errs.New(errs.Internal, err))
		return
	}

	appCustomer := toAppCustomer(bus)

	web.RenderResponse(http.StatusOK, w, r, appCustomer)
}

func (a *App) findByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var ctx = r.Context()

	parsedCustomerID, err := cvo.ParseCustomerID(id)
	if err != nil {
		web.RenderErrorResponse(ctx, w, r, err.Error(), err)
		return
	}

	bcus, err := a.srv.FindByIDAndOrgID(r.Context(), parsedCustomerID)
	if err != nil {
		web.RenderErrorResponse(ctx, w, r, err.Error(), errs.New(errs.NotFound, err))
		return
	}

	busids, err := a.identificationService.GetByCustomerID(ctx, bcus.ID.String())
	a.log.Info(ctx, "results from ids repo %v", busids)
	if err != nil {
		web.RenderErrorResponse(ctx, w, r, err.Error(), errs.New(errs.NotFound, err))
		return
	}

	appCustomer := toAppCustomer(bcus)
	appCustomer.WithIdentifications(&busids)
	web.RenderResponse(http.StatusOK, w, r, appCustomer)
}

func (a *App) upload(w http.ResponseWriter, r *http.Request) {
	var ctx = r.Context()
	var napp NewDocument

	if err := r.ParseMultipartForm(32 << 10); err != nil {
		a.log.Error(ctx, "parse form", err)
		return
	}

	napp.Classification = r.FormValue("classification")
	napp.DocumentType = r.FormValue("document_type")

	if err := napp.Validate(); err != nil {
		web.RenderErrorResponse(ctx, w, r, "validation failed", errs.New(errs.InvalidArgument, err))
		return
	}

	id := chi.URLParam(r, "id")
	parsedCustomerID, err := cvo.ParseCustomerID(id)
	if err != nil {
		web.RenderErrorResponse(ctx, w, r, err.Error(), err)
		return
	}

	bcus, err := a.srv.FindByIDAndOrgID(r.Context(), parsedCustomerID)
	if err != nil {
		web.RenderErrorResponse(ctx, w, r, err.Error(), errs.New(errs.NotFound, err))
		return
	}

	nbus, err := toBusNewDocument(napp)
	if err != nil {
		web.RenderErrorResponse(ctx, w, r, err.Error(), err)
		return
	}
	nbus.CustomerID = bcus.ID.String()

	a.log.Info(ctx, "bus for docus", nbus)

	file, header, err := r.FormFile("file")
	if err != nil {
		web.RenderErrorResponse(ctx, w, r, err.Error(), err)
		return
	}
	defer func() {
		_ = file.Close()
	}()

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		return
	}

	busFile, err := dvo.NewFile(header.Filename, header.Size, buf.Bytes())
	if err != nil {
		a.log.Error(ctx, "new file", err)
		web.RenderErrorResponse(ctx, w, r, err.Error(), err)
		return
	}
	nbus.File = busFile

	dbus, err := a.documentService.Create(r.Context(), nbus)
	if err != nil {
		a.log.Error(ctx, "document.service.create", err)

		web.RenderErrorResponse(ctx, w, r, err.Error(), err)
		return
	}
	appDoc := toAppDocument(dbus)
	web.RenderResponse(http.StatusCreated, w, r, appDoc)
}

func (a *App) createIdentifications(w http.ResponseWriter, r *http.Request) {
	var ctx = r.Context()
	id := chi.URLParam(r, "id")

	var napp NewIdentifications
	if err := json.NewDecoder(r.Body).Decode(&napp); err != nil {
		web.RenderErrorResponse(ctx, w, r, "invalid request",
			errs.New(errs.InvalidArgument, err))
		return
	}
	defer func() {
		err := r.Body.Close()
		if err != nil {
			return
		}
	}()

	if err := napp.Validate(); err != nil {
		web.RenderErrorResponse(ctx, w, r, "validation failed", errs.New(errs.InvalidArgument, err))
		return
	}

	parsedCustomerID, err := cvo.ParseCustomerID(id)
	if err != nil {
		web.RenderErrorResponse(ctx, w, r, err.Error(), err)
		return
	}

	bcus, err := a.srv.FindByIDAndOrgID(ctx, parsedCustomerID)
	if err != nil {
		web.RenderErrorResponse(ctx, w, r, err.Error(), errs.New(errs.NotFound, err))
		return
	}

	nbus, err := toBusNewIdentifications(napp, bcus)
	if err != nil {
		web.RenderErrorResponse(ctx, w, r, err.Error(), errs.New(errs.Internal, err))
		return
	}

	bus, err := a.identificationService.CreateMany(ctx, nbus)
	if err != nil {
		web.RenderErrorResponse(ctx, w, r, err.Error(), errs.New(errs.Internal, err))
		return
	}

	apps := make([]Identification, len(bus))
	for k, v := range bus {
		apps[k] = toAppIdentification(v)
	}

	web.RenderResponse(http.StatusCreated, w, r, apps)
}
