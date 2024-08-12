package rest

import (
	"encoding/json"
	"net/http"

	"github.com/afaridanquah/verifylab-backend/internal"
	"github.com/afaridanquah/verifylab-backend/internal/params"
	"github.com/afaridanquah/verifylab-backend/internal/service/verification"
	"github.com/go-chi/chi/v5"
)

type VerificationHandler struct {
	vs verification.VerificationService
}

func NewVerificationHandler(svc verification.VerificationService) *VerificationHandler {
	return &VerificationHandler{
		vs: svc,
	}
}

type Customer struct {
	ID          string  `json:"id"`
	FirstName   string  `json:"first_name"`
	LastName    string  `json:"last_name"`
	MiddleName  string  `json:"middle_name"`
	DateOfBirth string  `json:"date_of_birth"`
	Country     Country `json:"country"`
}

type GetVerificationResponse struct {
	ID               string   `json:"id"`
	VerificationType string   `json:"verification_type"`
	Customer         Customer `json:"customer"`
}

func (vh *VerificationHandler) Register(r chi.Router) {
	r.Post("/verifications", vh.create)
	// r.Get("/customers/{id}", ch.find)
}

func (vh *VerificationHandler) create(w http.ResponseWriter, r *http.Request) {
	var verificationReq params.CreateVerificationRequest

	if err := json.NewDecoder(r.Body).Decode(&verificationReq); err != nil {
		renderErrorResponse(w, r, "invalid request",
			internal.WrapErrorf(err, internal.ErrorCodeInvalidArgument, "json decoder"))
		return
	}

	defer r.Body.Close()

	ver, err := vh.vs.CreateVerification(r.Context(), verificationReq)

	if err != nil {
		renderErrorResponse(w, r, "create failed", err)
		return
	}

	renderResponse(w, r, ver.ID().String(), http.StatusCreated)

	// cus, err := ch.svc.CreateCustomer(r.Context(), customerReq)

}
