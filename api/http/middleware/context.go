package middleware

import (
	tilecontext "github.com/CXeon/tiles/context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func InjectContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader(tilecontext.HeaderTraceID) == "" {
			c.Request.Header.Set(tilecontext.HeaderTraceID, uuid.NewString())
		}

		appCtx := tilecontext.NewFromHTTPHeaders(c.Request.Context(), c.Request.Header)

		c.Request = c.Request.WithContext(appCtx)

		c.Next()
	}
}
