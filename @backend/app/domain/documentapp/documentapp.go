package documentapp

import (
	"encoding/json"
	"net/http"

	"bitbucket.org/msafaridanquah/verifylab-service/business/domain/documentbus"
	"bitbucket.org/msafaridanquah/verifylab-service/business/domain/documentbus/valueobject"
	"bitbucket.org/msafaridanquah/verifylab-service/foundation/ierr"
	"bitbucket.org/msafaridanquah/verifylab-service/foundation/logger"
	"bitbucket.org/msafaridanquah/verifylab-service/foundation/web"
)

type App struct {
	srv *documentbus.Service
	log *logger.Logger
}

func newApp(srv *documentbus.Service, log *logger.Logger) *App {
	return &App{
		srv: srv,
		log: log,
	}
}

func (app *App) upload(w http.ResponseWriter, r *http.Request) {
	var napp NewDocument
	if err := json.NewDecoder(r.Body).Decode(&napp); err != nil {
		web.RenderErrorResponse(app.log, w, r, "invalid request",
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
		web.RenderErrorResponse(app.log, w, r, "validation failed", ierr.WrapErrorf(err, ierr.ErrorCodeInvalidArgument, "json decoder"))
		return
	}

	nbus, err := toBusNewDocument(napp)
	if err != nil {
		web.RenderErrorResponse(app.log, w, r, err.Error(), err)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		web.RenderErrorResponse(app.log, w, r, err.Error(), err)
		return
	}

	busFile, err := valueobject.NewFile(header.Filename, header.Size, file)
	if err != nil {
		web.RenderErrorResponse(app.log, w, r, err.Error(), err)
		return
	}

	nbus.File = busFile

	bus, err := app.srv.Create(r.Context(), nbus)
	if err != nil {
		web.RenderErrorResponse(app.log, w, r, err.Error(), err)
		return
	}

	appDoc := toAppDocument(bus)

	web.RenderResponse(http.StatusCreated, w, r, appDoc)
}
