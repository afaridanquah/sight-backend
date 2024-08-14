package internal

import (
	"context"
	"fmt"
	"log"
	"net/url"

	"bitbucket.org/msafaridanquah/verifylab-service/internal"
	"bitbucket.org/msafaridanquah/verifylab-service/internal/envvar"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongoDB(conf *envvar.Configuration) (*mongo.Client, error) {
	get := func(v string) string {
		res, err := conf.Get(v)
		if err != nil {
			log.Fatalf("Couldn't get configuration value for %s: %s", v, err)
		}

		return res
	}

	databaseHost := get("MG_DATABASE_HOST")
	databasePort := get("MG_DATABASE_PORT")
	databaseUsername := get("MG_DATABASE_USERNAME")
	databasePassword := get("MG_DATABASE_PASSWORD")

	dsn := url.URL{
		Scheme: "mongodb",
		User:   url.UserPassword(databaseUsername, databasePassword),
		Host:   fmt.Sprintf("%s:%s", databaseHost, databasePort),
	}

	q := dsn.Query()

	dsn.RawQuery = q.Encode()

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dsn.String()))

	fmt.Printf("mongo string: %s", dsn.String())

	if err != nil {
		return nil, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "mongo.New")
	}

	return client, nil
}
