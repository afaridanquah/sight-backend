package rest

import (
	"encoding/json"
	"net/http"

	"github.com/afaridanquah/verifylab-service/internal"
	"github.com/afaridanquah/verifylab-service/internal/params"
	"github.com/afaridanquah/verifylab-service/internal/service/customer"
	"github.com/go-chi/chi/v5"
)

type CustomerHandler struct {
	svc customer.CustomerService
}

type FindCustomerResponse struct {
	ID          string  `json:"id"`
	FirstName   string  `json:"first_name"`
	LastName    string  `json:"last_name"`
	MiddleName  string  `json:"middle_name"`
	DateOfBirth string  `json:"date_of_birth"`
	Country     Country `json:"country"`
}

type Country struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

func (ch *CustomerHandler) Register(r chi.Router) {
	r.Post("/customers", ch.create)
	r.Get("/customers/{id}", ch.find)
}

func NewCustomerhandler(svc customer.CustomerService) *CustomerHandler {
	return &CustomerHandler{
		svc: svc,
	}
}

func (ch *CustomerHandler) create(w http.ResponseWriter, r *http.Request) {
	var customerReq params.CreateCustomerRequest
	if err := json.NewDecoder(r.Body).Decode(&customerReq); err != nil {
		renderErrorResponse(w, r, "invalid request",
			internal.WrapErrorf(err, internal.ErrorCodeInvalidArgument, "json decoder"))
		return
	}

	defer r.Body.Close()

	cus, err := ch.svc.CreateCustomer(r.Context(), customerReq)

	if err != nil {
		renderErrorResponse(w, r, "create failed", err)
		return
	}

	renderResponse(w, r, cus.ID().String(), http.StatusCreated)
}

func (ch *CustomerHandler) find(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	defer r.Body.Close()

	cus, err := ch.svc.FindCustomer(r.Context(), id)

	if err != nil {
		renderErrorResponse(w, r, "find failed", err)
		return
	}

	renderResponse(w, r, FindCustomerResponse{
		ID:          cus.ID().String(),
		FirstName:   cus.FirstName(),
		MiddleName:  cus.MiddleName(),
		LastName:    cus.Lastname(),
		DateOfBirth: string(cus.DateOfBirth()),
		Country: Country{
			Name: cus.Country().Name(),
			Code: cus.Country().Alpha2(),
		},
	}, http.StatusOK)

}
