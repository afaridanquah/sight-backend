package sdk

import (
	"os"

	"github.com/nats-io/nats.go"
)

func InitNatServer() (*nats.Conn, error) {
	natsURL := os.Getenv("NATS_URL")
	nc, err := nats.Connect(natsURL, nats.UserCredentials("NGS-Default-CLI.creds"), nats.Name("Sight"))
	if err != nil {
		return nil, err
	}

	// defer nc.Close()
	return nc, nil
}
