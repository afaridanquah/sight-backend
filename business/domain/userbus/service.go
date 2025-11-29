package userbus

import (
	"context"
	"log/slog"
	"time"

	"bitbucket.org/msafaridanquah/sight-backend/business/domain/userbus/valueobject"
	"bitbucket.org/msafaridanquah/sight-backend/foundation/logger"
	"bitbucket.org/msafaridanquah/sight-backend/foundation/otel"
	"github.com/mercari/go-circuitbreaker"
)

type Service struct {
	repo Repository
	cb   *circuitbreaker.CircuitBreaker
	log  *logger.Logger
}

type ServiceConfig func(*Service) error

func New(repo Repository, logger *logger.Logger, cfgs ...ServiceConfig) (*Service, error) {
	var ser = &Service{
		repo: repo,
		log:  logger,
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

func (s *Service) Create(ctx context.Context, nbus NewUser) (User, error) {
	ctx, span := otel.AddSpan(ctx, "tenantbus.service.create")
	defer span.End()

	now := time.Now()
	bus := User{
		ID:        valueobject.NewUserID(),
		FirstName: nbus.FirstName,
		LastName:  nbus.LastName,
		Email:     nbus.Email,
		TenantID:  nbus.TenantID,
		CreatedAt: now,
		UpdatedAt: now,
	}
	s.log.Info(ctx, "bus user", bus)
	if err := s.repo.Add(ctx, bus); err != nil {
		return User{}, err
	}

	return bus, nil
}

// func (s *Service) Authenticate(ctx context.Context, nbus LoginUserWithEmailAndPassword) (User, error) {
// 	ctx, span := otel.AddSpan(ctx, "tenantbus.service.authenticate")
// 	defer span.End()

// 	user, err := s.repo.GetByEmail(ctx, nbus.Email)
// 	s.log.Info(ctx, "user", user)

// 	if err != nil {
// 		return User{}, err
// 	}

// 	if !user.Password.Equals(nbus.Password) {
// 		s.log.Info(ctx, "password", nbus.Password)
// 		return User{}, errors.New("invalid credentials")
// 	}

// 	return user, nil
// }

// func (s *Service) QueryByUserID(ctx context.Context, id valueobject.ID) (User, error) {
// 	ctx, span := otel.AddSpan(ctx, "tenantbus.service.querybyuserid")
// 	defer span.End()

// 	user, err := s.repo.GetByID(ctx, id)
// 	if err != nil {
// 		return User{}, fmt.Errorf("query: userID[%s]  %w", id, err)
// 	}

// 	return user, nil
// }

// func (s *Service) Delete(ctx context.Context, id valueobject.ID) error {
// 	ctx, span := otel.AddSpan(ctx, "tenantbus.service.delete")
// 	defer span.End()

// 	if err := s.repo.Delete(ctx, id); err != nil {
// 		return fmt.Errorf("failed to execute delete action: %w", err)
// 	}

// 	return nil
// }

func (s *Service) Update(ctx context.Context, usr User, ups UpdateUser) (User, error) {
	ctx, span := otel.AddSpan(ctx, "tenantbus.service.update")
	defer span.End()

	if ups.Email != nil {
		usr.Email = *ups.Email
	}

	if ups.FirstName != nil {
		usr.FirstName = *ups.FirstName
	}

	if ups.LastName != nil {
		usr.LastName = *ups.LastName
	}

	if ups.OtherNames != nil {
		usr.OtherNames = *ups.OtherNames
	}

	if err := s.repo.Update(ctx, usr); err != nil {
		return User{}, err
	}

	return usr, nil
}

// func (s *Service) Query(ctx context.Context) (users []User, err error) {
// 	if !s.cb.Ready() {
// 		return []User{}, nil
// 	}

// 	defer func() {
// 		_ = s.cb.Done(ctx, err)
// 	}()

// 	users, err = s.repo.Query(ctx)
// 	if err != nil {
// 		return []User{}, err
// 	}

// 	return users, nil
// }
