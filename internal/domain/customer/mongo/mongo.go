package mongo

import (
	"context"
	"time"

	"github.com/afaridanquah/verifylab-service/internal/domain/customer"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoRepository struct {
	db        *mongo.Database
	customers *mongo.Collection
}

type mongoCustomer struct {
	ID        uuid.UUID `bson:"_id" json:"id"`
	FirstName string    `bson:"first_name" json:"first_name"`
	LastName  string    `bson:"last_name" json:"last_name"`
}

func NewFromCustomer(c customer.Customer) mongoCustomer {
	return mongoCustomer{
		ID:        c.GetID(),
		FirstName: c.GetFirstName(),
		LastName:  c.GetLastName(),
	}
}

func (m mongoCustomer) ToAggregate() customer.Customer {
	c := customer.Customer{}

	c.SetID(m.ID)
	c.SetFirstName(m.FirstName)
	c.SetLastName(m.LastName)

	return c

}
func New(ctx context.Context, connectionString string) (*MongoRepository, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionString))
	if err != nil {
		return nil, err
	}
	db := client.Database("verifylab")
	customers := db.Collection("customers")

	return &MongoRepository{
		db:        db,
		customers: customers,
	}, nil
}

func (mr *MongoRepository) GetByID(id uuid.UUID) (customer.Customer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result := mr.customers.FindOne(ctx, bson.M{"id": id})
	var c mongoCustomer
	err := result.Decode(&c)
	if err != nil {
		return customer.Customer{}, err
	}
	return c.ToAggregate(), nil
}

func (mr *MongoRepository) Add(cus customer.Customer) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	internal := NewFromCustomer(cus)
	_, err := mr.customers.InsertOne(ctx, internal)
	if err != nil {
		return err
	}
	return nil
}

func (mr *MongoRepository) Delete(id uuid.UUID) error {
	return nil
}

func (mr *MongoRepository) Update(cus customer.Customer) error {
	return nil
}
