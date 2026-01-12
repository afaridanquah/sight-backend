package documentbus

import (
	"context"
	_ "image/jpeg"
	_ "image/png"
	"log/slog"
	"time"

	"bitbucket.org/msafaridanquah/sight-backend/business/domain/documentbus/valueobject"
	"bitbucket.org/msafaridanquah/sight-backend/business/sdk/aws"
	"bitbucket.org/msafaridanquah/sight-backend/foundation/envvar"
	"bitbucket.org/msafaridanquah/sight-backend/foundation/logger"
	"bitbucket.org/msafaridanquah/sight-backend/foundation/otel"
	"github.com/mercari/go-circuitbreaker"
)

type Service struct {
	repo   Repository
	log    *logger.Logger
	cb     *circuitbreaker.CircuitBreaker
	envvar *envvar.Configuration
}

type ServiceConfig func(*Service) error

func New(repo Repository, log *logger.Logger, cfgs ...ServiceConfig) (*Service, error) {
	srv := &Service{
		log:  log,
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
			log.Info(context.Background(), "state changed",
				slog.String("old", string(oldState)),
				slog.String("new", string(newState)),
			)
		}),
	)

	return srv, nil
}

func (s *Service) WithEnv(envvar *envvar.Configuration) ServiceConfig {
	return func(s *Service) error {
		s.envvar = envvar
		return nil
	}
}

func (srv *Service) Create(ctx context.Context, nbus NewDocument) (Document, error) {
	ctx, span := otel.AddSpan(ctx, "documentbus.service.create")
	defer span.End()
	now := time.Now()

	status, err := valueobject.ParseStatus("PENDING")
	if err != nil {
		return Document{}, err
	}

	doc := Document{
		ID:             valueobject.NewDocumentID(),
		DocumentType:   nbus.DocumentType,
		Status:         status,
		Side:           nbus.Side,
		CustomerID:     nbus.CustomerID,
		BusinessID:     nbus.BusinessID,
		OriginalName:   nbus.File.OriginalName,
		FileName:       nbus.File.Name,
		Classification: nbus.Classification,
		CreatedAt:      now,
		UpdatedAt:      now,
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

func (srv *Service) UpdateStatus(ctx context.Context, doc Document, upd UpdateDocumentStatus) (Document, error) {
	ctx, span := otel.AddSpan(ctx, "documentbus.service.updatestatus")
	defer span.End()

	bus := Document{
		ID:             valueobject.NewDocumentID(),
		DocumentType:   doc.DocumentType,
		Status:         upd.Status,
		Side:           doc.Side,
		CustomerID:     doc.CustomerID,
		BusinessID:     doc.BusinessID,
		OriginalName:   doc.OriginalName,
		FileName:       doc.FileName,
		Classification: doc.Classification,
		CreatedAt:      doc.CreatedAt,
		UpdatedAt:      time.Now(),
	}

	if err := srv.repo.Update(ctx, bus); err != nil {
		return Document{}, nil
	}

	return bus, nil
}
