package command

import (
	"bitbucket.org/msafaridanquah/verifylab-service/foundation/logger"
	"bitbucket.org/msafaridanquah/verifylab-service/foundation/vaulti"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.opentelemetry.io/otel/trace"
)

type Config struct {
	Name       string
	Address    string
	PostgresDB *pgxpool.Pool
	Logger     *logger.Logger
	Tracer     trace.Tracer
	Vault      vaulti.Vaulty
}
