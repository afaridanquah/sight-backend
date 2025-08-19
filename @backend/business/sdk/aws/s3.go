package aws

import (
	"context"

	"bitbucket.org/msafaridanquah/verifylab-service/foundation/envvar"
	"bitbucket.org/msafaridanquah/verifylab-service/foundation/logger"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

const (
	AWS_S3_REGION = "us-east-2"  // Region
	AWS_S3_BUCKET = "sight.disk" // Bucket
)

type Config struct {
	Env *envvar.Configuration
	Log *logger.Logger
}

type AWSS3 struct {
	Client *s3.Client
}

func NewS3(conf Config) (*AWSS3, error) {
	get := func(v string) string {
		res, err := conf.Env.Get(v)
		if err != nil {
			conf.Log.Error(context.Background(), "env failed")
		}

		return res
	}

	key := get("AWS_ACCESS_KEY_ID")
	secret := get("AWS_SECRET_ACCESS_KEY")

	creds := credentials.NewStaticCredentialsProvider(key, secret, "")

	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithCredentialsProvider(creds))
	if err != nil {
		return &AWSS3{}, err
	}

	client := s3.NewFromConfig(cfg)

	return &AWSS3{
		Client: client,
	}, nil
}
