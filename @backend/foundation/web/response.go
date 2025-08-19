package web

import (
	"context"
	"errors"
	"net/http"

	"bitbucket.org/msafaridanquah/verifylab-service/foundation/ierr"
	"bitbucket.org/msafaridanquah/verifylab-service/foundation/logger"
	"github.com/go-chi/render"
)

// ErrorResponse represents a response containing an error message.
type ErrorResponse struct {
	Error       string           `json:"error"`
	Validations ierr.FieldErrors `json:"validations,omitempty"`
}

type httpStatus interface {
	HTTPStatus() int
}

func RenderErrorResponse(log *logger.Logger, w http.ResponseWriter, r *http.Request, msg string, err error) {
	resp := ErrorResponse{Error: msg}
	status := http.StatusInternalServerError

	// 	// If the context has been canceled, it means the client is no longer
	// 	// waiting for a response.
	if err := r.Context().Err(); err != nil {
		if errors.Is(err, context.Canceled) {
			err = errors.New("client disconnected, do not send response")
		}
	}

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
		_, span := addSpan(r.Context(), "rendererrorresponse")
		span.RecordError(err)
		defer span.End()
	}

	defer log.Error(r.Context(),
		"handled error during request",
		"err", err,
	)

	render.Status(r, status)
	render.JSON(w, r, &resp)
}

func RenderResponse(status int, w http.ResponseWriter, r *http.Request, data any) {
	render.Status(r, status)
	render.JSON(w, r, data)

}

// // Respond sends a response to the client.
// func Respond(ctx context.Context, w http.ResponseWriter, r *http.Request, data any, err error) error {
// 	// If the context has been canceled, it means the client is no longer
// 	// waiting for a response.
// 	if err := r.Context().Err(); err != nil {
// 		if errors.Is(err, context.Canceled) {
// 			err = errors.New("client disconnected, do not send response")
// 		}
// 	}

// 	statusCode := http.StatusOK

// 	switch v := resp.(type) {
// 	case httpStatus:
// 		statusCode = v.HTTPStatus()

// 	case error:
// 		statusCode = http.StatusInternalServerError

// 	default:
// 		if resp == nil {
// 			statusCode = http.StatusNoContent
// 		}
// 	}

// 	_, span := addSpan(ctx, "web.send.response", attribute.Int("status", statusCode))
// 	defer span.End()

// 	if statusCode == http.StatusNoContent {
// 		w.WriteHeader(statusCode)
// 		return nil
// 	}

// 	data, contentType, err := resp.Encode()
// 	if err != nil {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return fmt.Errorf("respond: encode: %w", err)
// 	}

// 	w.Header().Set("Content-Type", contentType)
// 	w.WriteHeader(statusCode)

// 	if _, err := w.Write(data); err != nil {
// 		return fmt.Errorf("respond: write: %w", err)
// 	}

// 	return nil
// }
