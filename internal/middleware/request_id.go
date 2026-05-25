package middleware

import (
	"context"
	"product-mall/internal/constants"
	"product-mall/pkg/pkg_logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func WithRequsetId() gin.HandlerFunc {
	return func(c *gin.Context) {
		uuID := c.GetHeader(constants.HeaderXRequestID)
		if len(uuID) == 0 {
			uuID = uuid.New().String()
		}

		ctx := context.WithValue(c.Request.Context(), constants.HeaderXRequestID, uuID)
		c.Request = c.Request.WithContext(ctx)
		c.Header(constants.HeaderXRequestID, uuID)

		if c.Request.Method == "GET" {
			pkg_logger.Logger.Info("request param", "query", c.Request.URL.Query())
		}

		c.Next()
	}
}