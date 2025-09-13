package otpapp

import (
	"bitbucket.org/msafaridanquah/sight-backend/business/domain/customerbus"
	"bitbucket.org/msafaridanquah/sight-backend/business/domain/otpbus"
	"bitbucket.org/msafaridanquah/sight-backend/foundation/logger"
	"github.com/go-chi/chi/v5"
)

type Config struct {
	Log             *logger.Logger
	Router          chi.Router
	OtpService      *otpbus.Service
	CustomerService *customerbus.Service
}

func Register(conf Config) {
	app := newApp(conf.OtpService, conf.CustomerService, conf.Log)

	conf.Router.Post("/customers/{id}/otps/new", app.create)
	conf.Router.Post("/customers/{id}/otps/verify", app.verify)
}
