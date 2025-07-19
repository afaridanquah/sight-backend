package customerbus

import (
	"context"
	"log/slog"
	"time"

	"bitbucket.org/msafaridanquah/verifylab-service/foundation/ierr"
	"bitbucket.org/msafaridanquah/verifylab-service/foundation/logger"
	"bitbucket.org/msafaridanquah/verifylab-service/foundation/otel"
	"github.com/google/uuid"
	"github.com/mercari/go-circuitbreaker"
)

type Service struct {
	customers Repository
	cb        *circuitbreaker.CircuitBreaker
}

type ServiceConfig func(*Service) error

func New(customers Repository, logger *logger.Logger, cfgs ...ServiceConfig) (*Service, error) {
	var ser = &Service{
		customers: customers,
	}

	for _, cfg := range cfgs {
		err := cfg(ser)
		if err != nil {
			return nil, err
		}
	}

	ser.cb = circuitbreaker.New(
		circuitbreaker.WithOpenTimeout(time.Minute*2),
		circuitbreaker.WithTripFunc(circuitbreaker.NewTripFuncConsecutiveFailures(3)),
		circuitbreaker.WithOnStateChangeHookFn(func(oldState, newState circuitbreaker.State) {
			logger.Info(context.Background(), "state changed",
				slog.String("old", string(oldState)),
				slog.String("new", string(newState)),
			)
		}),
	)

	return ser, nil
}

func WithRepository(cr Repository) ServiceConfig {
	return func(s *Service) error {
		s.customers = cr
		return nil
	}
}

func (cs *Service) Create(ctx context.Context, nc NewCustomer) (Customer, error) {
	ctx, span := otel.AddSpan(ctx, "customerbus.service.create")
	defer span.End()
	now := time.Now()

	userID := uuid.New()
	businessID := uuid.MustParse("6fe9cace-7c71-4e4b-b943-dd2f5bb21693")

	customer := Customer{
		ID:              uuid.New(),
		Person:          nc.Person,
		UserID:          userID,
		BusinessID:      businessID,
		DateOfBirth:     nc.DateOfBirth,
		Email:           nc.Email,
		PhoneNumber:     nc.PhoneNumber,
		BirthCountry:    nc.BirthCountry,
		Address:         nc.Address,
		Identifications: nc.Identifications,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	if err := cs.customers.Add(ctx, customer); err != nil {
		return Customer{}, ierr.WrapErrorf(err, ierr.ErrorCodeUnknown, "repo.create")
	}

	return customer, nil
}

func (cs *Service) QueryByIDAndBusinessID(ctx context.Context, custID uuid.UUID) (Customer, error) {
	ctx, span := otel.AddSpan(ctx, "service.customer.querybyidandbusinessid")
	defer span.End()

	//Get businessID from middleware
	businessID := uuid.MustParse("6fe9cace-7c71-4e4b-b943-dd2f5bb21693")

	cus, err := cs.customers.QueryByCustomerAndBusinessID(ctx, custID, businessID)
	if err != nil {
		return Customer{}, ierr.WrapErrorf(err, ierr.ErrorCodeUnknown, "repo.querybyidandbusinessid")
	}

	return cus, nil
}
