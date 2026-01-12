package natsio

import (
	"context"
	"encoding/json"
	"strings"

	"bitbucket.org/msafaridanquah/sight-backend/business/domain/customerbus"
	"bitbucket.org/msafaridanquah/sight-backend/foundation/logger"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type Repository struct {
	nc    *nats.Conn
	log   *logger.Logger
	topic string
}

func New(nats *nats.Conn, log *logger.Logger) *Repository {
	return &Repository{
		nc:    nats,
		log:   log,
		topic: "customers",
	}
}

func (r *Repository) Created(ctx context.Context, cus customerbus.Customer) error {
	if err := r.publish(ctx, cus); err != nil {
		return err
	}

	return nil
}

func (r *Repository) publish(ctx context.Context, cus customerbus.Customer) error {
	js, err := jetstream.New(r.nc)
	if err != nil {
		return err
	}

	_, err = js.CreateStream(ctx, jetstream.StreamConfig{
		Name:     r.topic,
		MaxBytes: 1024 * 1024 * 1024,
		Subjects: []string{"customers.*"},
	})

	if err != nil {
		return err
	}

	payload := Customer{
		ID:          cus.ID.String(),
		FirstName:   cus.Person.FirstName,
		LastName:    cus.Person.LastName,
		MiddleName:  cus.Person.MiddleName,
		DateOfBirth: cus.DateOfBirth.String(),
		Email:       cus.Email.String(),
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	_, err = js.Publish(ctx, strings.Join([]string{"customers", cus.OrgID}, "."), data)
	if err != nil {
		return err
	}
	return nil
}
