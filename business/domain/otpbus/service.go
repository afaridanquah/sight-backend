package otpbus

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"bitbucket.org/msafaridanquah/sight-backend/business/domain/otpbus/valueobject"
	"bitbucket.org/msafaridanquah/sight-backend/foundation/logger"
	"bitbucket.org/msafaridanquah/sight-backend/foundation/otel"
	"github.com/google/uuid"
	"github.com/mercari/go-circuitbreaker"
)

type ServiceConfig func(*Service) error

type Service struct {
	repo Repository
	cb   *circuitbreaker.CircuitBreaker
	log  *logger.Logger
}

func New(otps Repository, logger *logger.Logger, cfgs ...ServiceConfig) (*Service, error) {
	var ser = &Service{
		log:  logger,
		repo: otps,
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

func (srv *Service) Create(ctx context.Context, custID uuid.UUID, newbus NewOTP) (OTP, error) {
	ctx, span := otel.AddSpan(ctx, "otpbus.SendOtp")
	defer span.End()

	expiresAt := time.Now().Add(time.Minute * 10)

	code, err := generateOTPCode()
	if err != nil {
		return OTP{}, err
	}

	hashed, err := valueobject.ParseToHashCode(code)
	if err != nil {
		return OTP{}, err
	}

	bus := OTP{
		ID:          uuid.New(),
		Channel:     newbus.Channel,
		Destination: newbus.Destination,
		CustomerID:  custID,
		Code:        code,
		Hash:        hashed,
		ExpiresAt:   expiresAt,
	}

	defer srv.log.Info(ctx, "otp generated",
		slog.String("customer_id", bus.CustomerID.String()),
		slog.String("code", code),
	)

	if err := srv.repo.Add(ctx, bus); err != nil {
		return OTP{}, err
	}

	return bus, nil
}

func (srv *Service) Verify(ctx context.Context, newbus VerifyOTP) (OTP, error) {
	ctx, span := otel.AddSpan(ctx, "otpbus.VerifyOtpContent")
	defer span.End()

	hashed, err := valueobject.ParseToHashCode(newbus.Code)
	if err != nil {
		return OTP{}, fmt.Errorf("verify %w", err)
	}

	srv.log.Info(ctx, "verify.hashed", slog.String("hash", hashed.String()),
		slog.String("customerID", newbus.CustomerID.String()))

	res, err := srv.repo.FindByCustomerIDAndHash(ctx, newbus.CustomerID, hashed.String())

	if err != nil {
		return OTP{}, fmt.Errorf("verify %w", err)
	}

	if res.Code != newbus.Code {
		return OTP{}, fmt.Errorf("customer id %s :otp codes do not much", newbus.CustomerID)
	}

	// if err := srv.repo.Update(ctx, res.ID, res); err != nil {
	// 	return OTP{}, err
	// }

	return res, nil
}
