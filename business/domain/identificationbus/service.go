package identificationbus

import (
	"context"
	"log/slog"
	"time"

	"bitbucket.org/msafaridanquah/sight-backend/business/domain/identificationbus/valueobject"
	"bitbucket.org/msafaridanquah/sight-backend/foundation/logger"
	"bitbucket.org/msafaridanquah/sight-backend/foundation/otel"
	"github.com/mercari/go-circuitbreaker"
)

type Service struct {
	repo Repository
	log  *logger.Logger
	cb   *circuitbreaker.CircuitBreaker
}
type ServiceConfig func(*Service) error

func New(repo Repository, logger *logger.Logger, cfgs ...ServiceConfig) (*Service, error) {
	var ser = &Service{
		repo: repo,
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

func (srv *Service) Create(ctx context.Context, napp NewIdentification) (Identification, error) {
	ctx, span := otel.AddSpan(ctx, "identificationbus.service.create")
	defer span.End()

	now := time.Now()
	bus := Identification{
		ID:           valueobject.NewID(),
		CustomerID:   napp.CustomerID,
		FirstName:    napp.FirstName,
		LastName:     napp.LastName,
		MiddleName:   napp.MiddleName,
		Sex:          napp.Sex,
		Pin:          napp.Pin,
		PlaceOfBirth: napp.PlaceOfBirth,
		DateOfBirth:  napp.DateOfBirth,
		Nationality:  napp.Nationality,
		IssuedDate:   napp.IssuedDate,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := srv.repo.Add(ctx, bus); err != nil {
		return Identification{}, err
	}

	return bus, nil
}

func (srv *Service) CreateMany(ctx context.Context, napps []NewIdentification) ([]Identification, error) {
	ctx, span := otel.AddSpan(ctx, "identificationbus.service.create")
	defer span.End()

	now := time.Now()
	bus := make([]Identification, len(napps))

	if len(napps) > 0 {
		for k, v := range napps {
			bus[k] = Identification{
				ID:                 valueobject.NewID(),
				CustomerID:         v.CustomerID,
				FirstName:          v.FirstName,
				LastName:           v.LastName,
				MiddleName:         v.MiddleName,
				Sex:                v.Sex,
				Pin:                v.Pin,
				IdentificationType: v.IdentificationType,
				PlaceOfBirth:       v.PlaceOfBirth,
				DateOfBirth:        v.DateOfBirth,
				Nationality:        v.Nationality,
				IssuedDate:         v.IssuedDate,
				ExpDate:            v.ExpDate,
				Country:            v.Country,
				CreatedAt:          now,
				UpdatedAt:          now,
			}
		}
	}

	if err := srv.repo.AddMany(ctx, bus); err != nil {
		return []Identification{}, err
	}

	return bus, nil
}

func (srv *Service) GetByCustomerID(ctx context.Context, cis string) ([]Identification, error) {
	ctx, span := otel.AddSpan(ctx, "identificationbus.service.findbycustomerid")
	defer span.End()

	ci, err := valueobject.ParseCustomerID(cis)
	if err != nil {
		return []Identification{}, err
	}

	res, err := srv.repo.QueryByCustomerID(ctx, ci)
	if err != nil {
		return []Identification{}, err
	}

	return res, nil
}
