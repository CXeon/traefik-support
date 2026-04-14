package forwardauth

import (
	appforwardauth "github.com/CXeon/traefik-support/internal/application/forwardauth"
	forwardauthhttp "github.com/CXeon/traefik-support/internal/api/http/forwardauth"
	forwardauthdomain "github.com/CXeon/traefik-support/internal/domain/forwardauth"
	infraforwardauth "github.com/CXeon/traefik-support/internal/infrastructure/forwardauth"
	"github.com/CXeon/traefik-support/internal/modular"
	"github.com/gin-gonic/gin"
)

type Module struct{}

func New() *Module { return &Module{} }

func (m *Module) Build(_ modular.Deps, g *gin.RouterGroup) error {
	userClient := infraforwardauth.NewStubUserClient()
	domainSvc := forwardauthdomain.NewAuthDomainService(userClient)
	uc := appforwardauth.NewUseCase(domainSvc)
	handler := forwardauthhttp.NewHandler(uc)
	handler.Register(g)
	return nil
}
