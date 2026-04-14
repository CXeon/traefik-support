package modular

import (
	tilesLogger "github.com/CXeon/tiles/logger"
	infraRegistry "github.com/CXeon/traefik-support/internal/infrastructure/registry"
	"github.com/gin-gonic/gin"
)

type Deps struct {
	Logger         tilesLogger.Logger
	ServiceLocator infraRegistry.Locator // nil when registry disabled
}

type Module interface {
	Build(deps Deps, g *gin.RouterGroup) error
}
