package forwardauth

import (
	"errors"
	"net/http"
	"strings"

	tilecontext "github.com/CXeon/tiles/context"
	tileerrors "github.com/CXeon/tiles/errors"
	traefikheader "github.com/CXeon/tiles/gateway/traefik"
	"github.com/CXeon/traefik-support/api/http/response"
	appforwardauth "github.com/CXeon/traefik-support/internal/application/forwardauth"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	useCase appforwardauth.UseCase
}

func NewHandler(useCase appforwardauth.UseCase) *Handler {
	return &Handler{useCase: useCase}
}

func (h *Handler) Register(g *gin.RouterGroup) {
	// Traefik ForwardAuth forwards the original request method to the auth endpoint,
	// so we must accept any method, not just POST.
	g.Any("/forward-auth", h.Handle)
}

// Handle 是 Traefik ForwardAuth 中间件调用的认证接口。
// 认证成功返回 200（Traefik 放行），失败必须返回 401（不能用通用 response.Fail，它返回 200）。
func (h *Handler) Handle(c *gin.Context) {
	req := &appforwardauth.CheckRequest{
		Env:     c.GetHeader(traefikheader.HeaderKeyEnv),
		Cluster: c.GetHeader(traefikheader.HeaderKeyCluster),
		Company: c.GetHeader(traefikheader.HeaderKeyCompany),
		Project: c.GetHeader(traefikheader.HeaderKeyProject),
		Color:   c.GetHeader(traefikheader.HeaderKeyColor),
		UserID:  c.GetHeader("X-UserID"),
		JWT:     extractBearer(c),
		Path:    c.GetHeader("X-Forwarded-Uri"),
	}

	err := h.useCase.Check(c.Request.Context(), req)
	if err != nil {
		traceID := tilecontext.From(c.Request.Context()).TraceID()
		var tileErr *tileerrors.Error
		if errors.As(err, &tileErr) {
			c.JSON(http.StatusUnauthorized, response.Response{
				Code:    tileErr.Code,
				Message: tileErr.Message,
				TraceID: traceID,
			})
		} else {
			c.JSON(http.StatusUnauthorized, response.Response{
				Code:    tileerrors.ErrUnauthorized.Code,
				Message: tileerrors.ErrUnauthorized.Message,
				TraceID: traceID,
			})
		}
		return
	}

	c.Status(http.StatusOK)
}

func extractBearer(c *gin.Context) string {
	header := c.GetHeader("Authorization")
	if !strings.HasPrefix(header, "Bearer ") {
		return ""
	}
	return strings.TrimPrefix(header, "Bearer ")
}
