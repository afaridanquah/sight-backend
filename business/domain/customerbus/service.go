package customerbus

import (
	"context"
	"log/slog"
	"time"

	"bitbucket.org/msafaridanquah/sight-backend/business/domain/customerbus/valueobject"
	"bitbucket.org/msafaridanquah/sight-backend/foundation/logger"
	"bitbucket.org/msafaridanquah/sight-backend/foundation/otel"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mercari/go-circuitbreaker"
)

type Service struct {
	repo Repository
	cb   *circuitbreaker.CircuitBreaker
	es   SearchRepository
	msg  MessageBrokerRepository
	pool *pgxpool.Pool
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

func WithMsgBroker(msg MessageBrokerRepository) ServiceConfig {
	return func(s *Service) error {
		s.msg = msg
		return nil
	}
}

func (srv *Service) Create(ctx context.Context, nc NewCustomer) (Customer, error) {
	ctx, span := otel.AddSpan(ctx, "customerbus.service.create")
	defer span.End()
	now := time.Now()

	orgID, err := valueobject.ParseOrgID("org_372gfpgRzfQnwHDhlVAH00orIUR")
	if err != nil {
		return Customer{}, err
	}

	customer := Customer{
		ID:     valueobject.NewCustomerID(),
		Person: nc.Person,
		// UserID:          nc.UserID,
		OrgID:           orgID.String(),
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

	// if err := srv.msg.Created(ctx, customer); err != nil {
	// 	return Customer{}, err
	// }

	return customer, nil
}

func (srv *Service) Update(ctx context.Context, cust Customer, up UpdateCustomer) (Customer, error) {
	ctx, span := otel.AddSpan(ctx, "customerbus.service.update")
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

func (srv *Service) FindByIDAndOrgID(ctx context.Context, custID valueobject.ID) (Customer, error) {
	ctx, span := otel.AddSpan(ctx, "customerbus.service.querybyidandbusinessid")
	defer span.End()

	// Get businessID from middleware
	orgID, err := valueobject.ParseOrgID("org_372gfpgRzfQnwHDhlVAH00orIUR")
	if err != nil {
		return Customer{}, err
	}

	cus, err := srv.repo.QueryByIDAndOrgID(ctx, custID, orgID)
	if err != nil {
		return Customer{}, err
	}

	return cus, nil
}

func (srv *Service) QueryByOrgID(ctx context.Context) ([]Customer, error) {
	ctx, span := otel.AddSpan(ctx, "customerbus.service.query")
	defer span.End()

	orgID, err := valueobject.ParseOrgID("org_372gfpgRzfQnwHDhlVAH00orIUR")
	if err != nil {
		return []Customer{}, err
	}

	customers, err := srv.repo.QueryByOrgID(ctx, orgID)
	if err != nil {
		return []Customer{}, err
	}

	return customers, nil
}
