package otpbus

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"bitbucket.org/msafaridanquah/verifylab-service/foundation/logger"
	"bitbucket.org/msafaridanquah/verifylab-service/foundation/otel"
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

	code, hashed, err := generateOTPCode()
	if err != nil {
		return OTP{}, err
	}

	bus := OTP{
		ID:          uuid.New(),
		Channel:     newbus.Channel,
		Destination: newbus.Destination,
		CustomerID:  custID,
		HashedCode:  hashed,
		ExpiresAt:   expiresAt,
	}

	srv.log.Info(ctx, "otp generated",
		slog.String("customer_id", bus.CustomerID.String()),
		slog.String("customer_id", code),
		slog.String("customer_id", hashed.String()),
	)

	res, err := srv.repo.Add(ctx, bus)
	if err != nil {
		return OTP{}, err
	}

	return res, nil
}

func (srv *Service) Verify(ctx context.Context, custID uuid.UUID, newbus VerifyOTP) (OTP, error) {
	ctx, span := otel.AddSpan(ctx, "otpbus.VerifyOtpContent")
	defer span.End()

	res, err := srv.repo.FindByCustomerIDAndHash(ctx, custID, newbus.HashedCode)

	if err != nil {
		return OTP{}, fmt.Errorf("verify %w", err)
	}

	if res.HashedCode.NotEqual(newbus.HashedCode) {
		return OTP{}, fmt.Errorf("customer id %s :otp codes do not much", custID)
	}

	if err := srv.repo.Update(ctx, res.ID, res); err != nil {
		return OTP{}, err
	}

	return res, nil
}
