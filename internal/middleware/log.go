package middleware

import (
	"product-mall/pkg/pkg_logger"
	"time"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		c.Next()
		endTime := time.Now()
		latencyTime := endTime.Sub(startTime)

		pkg_logger.LogrusObj.Info("request completed",
			"status", c.Writer.Status(),
			"latency", latencyTime,
			"client_ip", c.ClientIP(),
			"method", c.Request.Method,
			"uri", c.Request.RequestURI,
			"user_agent", c.Request.UserAgent(),
			"referer", c.Request.Referer(),
		)
	}
}