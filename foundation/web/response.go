package web

import (
	"context"
	"errors"
	"net/http"

	"bitbucket.org/msafaridanquah/verifylab-service/foundation/ierr"
	"github.com/go-chi/render"
)

// ErrorResponse represents a response containing an error message.
type ErrorResponse struct {
	Error       string           `json:"error"`
	Validations ierr.FieldErrors `json:"validations,omitempty"`
}

// =============================================================================

func RenderErrorResponse(ctx context.Context, w http.ResponseWriter, r *http.Request, msg string, err error) {
	resp := ErrorResponse{Error: msg}
	status := http.StatusInternalServerError

	// // 	// If the context has been canceled, it means the client is no longer
	// // 	// waiting for a response.
	// if err := ctx.Err(); err != nil {
	// 	if errors.Is(err, context.Canceled) {
	// 		err = errors.New("client disconnected, do not send response")
	// 	}
	// }

	var aerr *ierr.Error

	if !errors.As(err, &aerr) {
		resp.Error = "internal error"
	} else {
		switch aerr.Code() {
		case ierr.InvalidArgument:
			status = http.StatusBadRequest

			var verrors ierr.FieldErrors

			if errors.As(aerr, &verrors) {
				resp.Validations = verrors
			}
		case ierr.Unknown:
			fallthrough
		default:
			status = aerr.HTTPStatus()
		}
	}

	if err != nil {
		_, span := addSpan(ctx, "rendererrorresponse")

		span.RecordError(err)
		defer span.End()
	}

	render.Status(r, status)
	render.JSON(w, r, &resp)
}

func RenderResponse(status int, w http.ResponseWriter, r *http.Request, data any) {
	render.Status(r, status)
	render.JSON(w, r, data)
}
