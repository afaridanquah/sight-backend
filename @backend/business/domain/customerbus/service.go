package customerbus

import (
	"context"
	"log/slog"
	"time"

	"bitbucket.org/msafaridanquah/verifylab-service/business/domain/customerbus/valueobject"
	"bitbucket.org/msafaridanquah/verifylab-service/foundation/logger"
	"bitbucket.org/msafaridanquah/verifylab-service/foundation/otel"
	"github.com/google/uuid"
	"github.com/mercari/go-circuitbreaker"
)

type Service struct {
	repo Repository
	cb   *circuitbreaker.CircuitBreaker
}

type ServiceConfig func(*Service) error

func New(repo Repository, logger *logger.Logger, cfgs ...ServiceConfig) (*Service, error) {
	var srv = &Service{
		repo: repo,
	}

	for _, cfg := range cfgs {
		err := cfg(srv)
		if err != nil {
			return nil, err
		}
	}

	srv.cb = circuitbreaker.New(
		circuitbreaker.WithOpenTimeout(time.Minute*2),
		circuitbreaker.WithTripFunc(circuitbreaker.NewTripFuncConsecutiveFailures(3)),
		circuitbreaker.WithOnStateChangeHookFn(func(oldState, newState circuitbreaker.State) {
			logger.Info(context.Background(), "state changed",
				slog.String("old", string(oldState)),
				slog.String("new", string(newState)),
			)
		}),
	)

	return srv, nil
}

func WithRepository(cr Repository) ServiceConfig {
	return func(s *Service) error {
		s.repo = cr
		return nil
	}
}

func (srv *Service) Create(ctx context.Context, nc NewCustomer) (Customer, error) {
	ctx, span := otel.AddSpan(ctx, "customerbus.service.create")
	defer span.End()
	now := time.Now()

	userID := uuid.New()
	orgID := uuid.MustParse("6fe9cace-7c71-4e4b-b943-dd2f5bb21693")

	customer := Customer{
		ID:              uuid.New(),
		Person:          nc.Person,
		UserID:          userID,
		OrgID:           orgID,
		DateOfBirth:     nc.DateOfBirth,
		Email:           nc.Email,
		PhoneNumber:     nc.PhoneNumber,
		BirthCountry:    nc.BirthCountry,
		Address:         nc.Address,
		Identifications: nc.Identifications,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	if err := srv.repo.Add(ctx, customer); err != nil {
		return Customer{}, err
	}

	return customer, nil
}

func (srv *Service) Update(ctx context.Context, cust Customer, up UpdateCustomer) (Customer, error) {
	ctx, span := otel.AddSpan(ctx, "business.customerbus.service.update")
	defer span.End()

	if up.Person != nil {
		cust.Person = *up.Person
	}

	if up.BirthCountry != nil {
		cust.BirthCountry = *up.BirthCountry
	}

	if up.Email != nil {
		cust.Email = *up.Email
	}

	if up.Address != nil {
		cust.Address = *up.Address
	}

	if up.Identifications != nil {
		idx := make([]valueobject.Identification, len(*up.Identifications))
		copy(idx, *up.Identifications)
		cust.Identifications = idx
	}

	if err := srv.repo.Update(ctx, cust); err != nil {
		return Customer{}, err
	}

	return cust, nil
}

func (srv *Service) FindByIDAndOrgID(ctx context.Context, custID uuid.UUID) (Customer, error) {
	ctx, span := otel.AddSpan(ctx, "business.customerbus.service.querybyidandbusinessid")
	defer span.End()

	//Get businessID from middleware
	orgID := uuid.MustParse("6fe9cace-7c71-4e4b-b943-dd2f5bb21693")

	cus, err := srv.repo.QueryByIDAndOrgID(ctx, custID, orgID)
	if err != nil {
		return Customer{}, err
	}

	return cus, nil
}
