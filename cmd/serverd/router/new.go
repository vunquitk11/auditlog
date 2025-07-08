package router

import (
	"context"

	"github.com/audit-log-service/internal/controller/auditlog"
	"github.com/audit-log-service/internal/handler/rest/authenticated"
	"github.com/audit-log-service/internal/repository"
)

type Router struct {
	authenticatedRESTHandler *authenticated.Handler
}

func New(
	ctx context.Context,
	corsOrigins []string,
	repo repository.Registry,
) Router {
	auditCtrl := auditlog.New(repo)
	authHandler := authenticated.New(auditCtrl)
	return Router{
		authenticatedRESTHandler: &authHandler,
	}
}
