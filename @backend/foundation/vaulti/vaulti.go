package vaulti

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"os"

	"bitbucket.org/msafaridanquah/verifylab-service/foundation/logger"
	"github.com/hashicorp/vault-client-go"
	"github.com/hashicorp/vault-client-go/schema"
)

type Config struct {
	KeyName string
	Log     *logger.Logger
}

type Vaulty struct {
	client *vault.Client
	log    *logger.Logger
}

type TransitResponse struct {
	// Specify the encrypted payload
	Ciphertext string `json:"ciphertext,omitempty"`

	Plaintext  string         `json:"plaintext,omitempty"`
	BulkData   map[string]any `json:"bulk,omitempty"`
	KeyVersion int            `json:"key_version,omitempty"`
}

type TransitInput struct {
	Ciphertext string `json:"ciphertext"`
}

type TransitBatchInput struct {
	BatchInput []TransitInput `json:"batch_input"`
}

type TransitDecryptResponse struct {
	Plaintext string `json:"plaintext"`
}

func InitVault(conf Config) (*Vaulty, error) {
	host := "http://vault:8300"
	vaultToken := os.Getenv("VAULT_TOKEN")

	client, err := vault.New(
		vault.WithAddress(host),
	)

	client.SetToken(vaultToken)
	if err != nil {
		return &Vaulty{}, err
	}

	return &Vaulty{
		client: client,
		log:    conf.Log,
	}, nil
}

func (v *Vaulty) TransitEncrypt(data string, key string) (TransitResponse, error) {
	tob64 := base64.StdEncoding.EncodeToString([]byte(data))
	res, err := v.client.Secrets.TransitEncrypt(context.Background(), key, schema.TransitEncryptRequest{
		Plaintext: tob64,
	})
	if err != nil {
		return TransitResponse{}, err
	}

	var transitResponse TransitResponse

	bytes, err := json.Marshal(res.Data)
	if err != nil {
		return TransitResponse{}, err
	}

	if err := json.Unmarshal(bytes, &transitResponse); err != nil {
		return TransitResponse{}, err
	}

	return transitResponse, nil
}

func (v *Vaulty) TransitDecrypt(ciphertext string, key string) (TransitResponse, error) {
	res, err := v.client.Secrets.TransitDecrypt(context.Background(), key, schema.TransitDecryptRequest{
		Ciphertext: ciphertext,
	})
	if err != nil {
		return TransitResponse{}, err
	}

	var transitResponse TransitResponse

	bytes, err := json.Marshal(res.Data)
	if err != nil {
		return TransitResponse{}, err
	}

	if err := json.Unmarshal(bytes, &transitResponse); err != nil {
		return TransitResponse{}, err
	}

	decoded, err := base64.StdEncoding.DecodeString(string(transitResponse.Plaintext))
	if err := json.Unmarshal(bytes, &transitResponse); err != nil {
		return TransitResponse{}, err
	}

	transitResponse.Plaintext = string(decoded)

	return transitResponse, nil
}
