package organizationapi

import (
	"bitbucket.org/msafaridanquah/sight-backend/app/domain/organizationapp"
	"bitbucket.org/msafaridanquah/sight-backend/business/domain/organizationbus"
	"bitbucket.org/msafaridanquah/sight-backend/business/domain/organizationbus/postgres"
	"bitbucket.org/msafaridanquah/sight-backend/foundation/envvar"
	"bitbucket.org/msafaridanquah/sight-backend/foundation/logger"
	"bitbucket.org/msafaridanquah/sight-backend/foundation/vaulti"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Routes(logger *logger.Logger, pool *pgxpool.Pool, envvar *envvar.Configuration, vaulti *vaulti.Vaulty, chi chi.Router) {
	repo := postgres.New(pool, pool, vaulti)

	service, _ := organizationbus.New(repo, logger)

	organizationapp.Register(organizationapp.Config{
		Log:        logger,
		OrgService: service,
		Router:     chi,
		Repo:       repo,
	})
}
