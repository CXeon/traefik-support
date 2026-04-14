package router

import (
	"fmt"

	"github.com/CXeon/traefik-support/api/http/middleware"
	"github.com/CXeon/traefik-support/internal/config"
	"github.com/gin-gonic/gin"
)

func Default(r *gin.Engine) *gin.RouterGroup {
	prefix := fmt.Sprintf("/%s/%s/%s/api/v1",
		config.Config.Company,
		config.Config.Project,
		config.Config.ServiceName,
	)
	r.Use(middleware.CORS(), gin.Recovery(), middleware.InjectContext())
	g := r.Group(prefix)
	g.GET("/health", func(c *gin.Context) {
		c.String(200, "OK")
	})
	return g
}
