package customerapp

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	"bitbucket.org/msafaridanquah/verifylab-service/business/domain/customerbus"
	"bitbucket.org/msafaridanquah/verifylab-service/business/domain/documentbus"
	dvo "bitbucket.org/msafaridanquah/verifylab-service/business/domain/documentbus/valueobject"
	"bitbucket.org/msafaridanquah/verifylab-service/foundation/envvar"
	"bitbucket.org/msafaridanquah/verifylab-service/foundation/ierr"
	"bitbucket.org/msafaridanquah/verifylab-service/foundation/logger"
	"bitbucket.org/msafaridanquah/verifylab-service/foundation/web"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type App struct {
	srv             *customerbus.Service
	documentService *documentbus.Service
	log             *logger.Logger
	envvar          *envvar.Configuration
}

func newApp(srv *customerbus.Service, ds *documentbus.Service, log *logger.Logger, envvar *envvar.Configuration) *App {
	return &App{
		srv:             srv,
		log:             log,
		documentService: ds,
		envvar:          envvar,
	}
}

func (a *App) create(w http.ResponseWriter, r *http.Request) {
	var napp NewCustomer
	var ctx = r.Context()
	if err := json.NewDecoder(r.Body).Decode(&napp); err != nil {
		a.log.Error(ctx, "documentapp.create", slog.String("err", err.Error()))

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
		a.log.Error(ctx, "customerapp.validate", err)
		web.RenderErrorResponse(ctx, w, r, "validation failed", ierr.WrapErrorf(err, ierr.InvalidArgument, "json decoder"))
		return
	}

	newbus, err := toBusNewCustomer(napp)
	a.log.Info(ctx, "customerapp.toBusNewCustomer", newbus)
	if err != nil {
		web.RenderErrorResponse(ctx, w, r, err.Error(), err)
		return
	}

	bcus, err := a.srv.Create(r.Context(), newbus)
	if err != nil {
		a.log.Info(ctx, "customerapp.srv.create", err)
		web.RenderErrorResponse(ctx, w, r, err.Error(), err)
		return
	}

	appCustomer := toAppCustomer(bcus)
	web.RenderResponse(http.StatusCreated, w, r, appCustomer)
}

func (a *App) findByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var ctx = r.Context()

	parsedCustomerID, err := uuid.Parse(id)
	if err != nil {
		web.RenderErrorResponse(ctx, w, r, err.Error(), err)
		return
	}

	bcus, err := a.srv.FindByIDAndOrgID(r.Context(), parsedCustomerID)
	if err != nil {
		web.RenderErrorResponse(ctx, w, r, err.Error(), ierr.WrapErrorf(err, ierr.NotFound, "querybyidandbusinessid"))
		return
	}

	appCustomer := toAppCustomer(bcus)
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
		web.RenderErrorResponse(ctx, w, r, "validation failed", ierr.WrapErrorf(err, ierr.InvalidArgument, "json decoder"))
		return
	}

	id := chi.URLParam(r, "id")
	parsedCustomerID, err := uuid.Parse(id)
	if err != nil {
		web.RenderErrorResponse(ctx, w, r, err.Error(), err)
		return
	}

	bcus, err := a.srv.FindByIDAndOrgID(r.Context(), parsedCustomerID)
	if err != nil {
		web.RenderErrorResponse(ctx, w, r, err.Error(), ierr.WrapErrorf(err, ierr.NotFound, "querybyidandbusinessid"))
		return
	}

	nbus, err := toBusNewDocument(napp)
	if err != nil {
		web.RenderErrorResponse(ctx, w, r, err.Error(), err)
		return
	}
	nbus.CustomerID = bcus.ID

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
