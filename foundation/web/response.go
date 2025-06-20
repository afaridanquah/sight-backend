package web

import (
	"errors"
	"net/http"

	"bitbucket.org/msafaridanquah/verifylab-service/foundation/ierr"
	"bitbucket.org/msafaridanquah/verifylab-service/foundation/otel"
	"github.com/go-chi/render"
)

// ErrorResponse represents a response containing an error message.
type ErrorResponse struct {
	Error       string           `json:"error"`
	Validations ierr.FieldErrors `json:"validations,omitempty"`
}

func RenderErrorResponse(w http.ResponseWriter, r *http.Request, msg string, err error) {
	resp := ErrorResponse{Error: msg}
	status := http.StatusInternalServerError

	var aerr *ierr.Error

	if !errors.As(err, &aerr) {
		resp.Error = "internal error"
	} else {

		switch aerr.Code() {
		case ierr.ErrorCodeNotFound:
			status = http.StatusNotFound
		case ierr.ErrorCodeInvalidArgument:
			status = http.StatusBadRequest

			var verrors ierr.FieldErrors

			if errors.As(aerr, &verrors) {
				resp.Validations = verrors
			}
		case ierr.ErrorCodeUnknown:
			fallthrough
		default:
			status = http.StatusInternalServerError
		}
	}

	if err != nil {
		// _, span := otel.Tracer().Start(r.Context(), "renderErrorResponse")
		_, span := otel.AddSpan(r.Context(), "renderErrorResponse")
		span.RecordError(err)
		defer span.End()
	}

	render.Status(r, status)
	render.JSON(w, r, &resp)
}

func RenderResponse(w http.ResponseWriter, r *http.Request, res interface{}, status int) {
	render.Status(r, status)
	render.JSON(w, r, res)
}
