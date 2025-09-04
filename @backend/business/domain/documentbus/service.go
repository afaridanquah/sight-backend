package documentbus

import (
	"context"
	"log/slog"
	"time"

	"bitbucket.org/msafaridanquah/verifylab-service/business/sdk/aws"
	"bitbucket.org/msafaridanquah/verifylab-service/foundation/envvar"
	"bitbucket.org/msafaridanquah/verifylab-service/foundation/logger"
	"bitbucket.org/msafaridanquah/verifylab-service/foundation/otel"
	"github.com/google/uuid"
	"github.com/mercari/go-circuitbreaker"
)

type Service struct {
	repo   Repository
	log    *logger.Logger
	cb     *circuitbreaker.CircuitBreaker
	envvar *envvar.Configuration
}

type ServiceConfig func(*Service) error

func New(log *logger.Logger, cfgs ...ServiceConfig) (*Service, error) {
	srv := &Service{
		log: log,
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
			log.Info(context.Background(), "state changed",
				slog.String("old", string(oldState)),
				slog.String("new", string(newState)),
			)
		}),
	)

	return srv, nil
}

func WithRepository(repo Repository) ServiceConfig {
	return func(s *Service) error {
		s.repo = repo
		return nil
	}
}

func WithEnv(envvar *envvar.Configuration) ServiceConfig {
	return func(s *Service) error {
		s.envvar = envvar
		return nil
	}
}

func (srv *Service) Create(ctx context.Context, nbus NewDocument) (Document, error) {
	ctx, span := otel.AddSpan(ctx, "documentbus.service.create")
	defer span.End()
	now := time.Now()

	doc := Document{
		ID:             uuid.New(),
		DocumentType:   nbus.DocumentType,
		Side:           nbus.Side,
		CustomerID:     nbus.CustomerID,
		BusinessID:     nbus.BusinessID,
		OriginalName:   nbus.File.OriginalName,
		FileName:       nbus.File.Name,
		CreatedAt:      now,
		UpdatedAt:      now,
		Classification: nbus.Classification,
	}

	s3, err := aws.NewS3(aws.Config{
		Env: srv.envvar,
		Log: srv.log,
	})
	if err != nil {
		return Document{}, err
	}

	if err := s3.Upload(nbus.File.Data, nbus.File.Name); err != nil {
		return Document{}, err
	}

	if err := srv.repo.Add(ctx, doc); err != nil {
		return Document{}, err
	}

	return doc, nil
}
