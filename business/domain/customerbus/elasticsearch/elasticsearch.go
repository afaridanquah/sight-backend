package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"bitbucket.org/msafaridanquah/sight-backend/business/domain/customerbus"
	es "github.com/elastic/go-elasticsearch/v9"
	esapi "github.com/elastic/go-elasticsearch/v9/esapi"
)

type Repository struct {
	client *es.Client
	index  string
}

func New(client *es.Client) *Repository {
	return &Repository{
		client: client,
		index:  "customers",
	}
}

func (r *Repository) Index(ctx context.Context, customerbus customerbus.Customer) error {
	body := Customer{
		ID:          customerbus.ID.String(),
		FirstName:   customerbus.Person.FirstName,
		LastName:    customerbus.Person.LastName,
		MiddleName:  customerbus.Person.MiddleName,
		DateOfBirth: customerbus.DateOfBirth.String(),
		Email:       customerbus.Email.String(),
		CreatedAt:   customerbus.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   customerbus.UpdatedAt.Format(time.RFC3339),
	}

	var buf bytes.Buffer

	if err := json.NewEncoder(&buf).Encode(body); err != nil {
		return fmt.Errorf("encode json %w:", err)
	}

	req := esapi.IndexRequest{
		Index:      r.index,
		Body:       &buf,
		DocumentID: customerbus.ID.String(),
		Refresh:    "true",
	}

	res, err := req.Do(context.Background(), r.client)
	if err != nil {
		return err
	}

	defer func() {
		_ = res.Body.Close()
	}()

	return nil
}
