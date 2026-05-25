package middleware

import (
	"net/http"
	util "product-mall/internal/tools"
	"product-mall/pkg/e"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var code int
		code = e.SUCCESS

		token := extractToken(c)
		if token == "" {
			code = e.ErrorAuthCheckTokenFail
		} else {
			claims, err := util.ParseToken(token)
			if err != nil {
				code = e.ErrorAuthCheckTokenFail
			} else if claims.ExpiresAt.Before(time.Now()) {
				code = e.ErrorAuthCheckTokenTimeout
			} else {
				c.Set("user_id", claims.ID)
				c.Set("username", claims.Username)
			}
		}

		if code != e.SUCCESS {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status": code,
				"msg":    e.GetMsg(code),
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

func JWTAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		var code int
		code = e.SUCCESS

		token := extractToken(c)
		if token == "" {
			code = e.ErrorAuthCheckTokenFail
		} else {
			claims, err := util.ParseToken(token)
			if err != nil {
				code = e.ErrorAuthCheckTokenFail
			} else if claims.ExpiresAt.Before(time.Now()) {
				code = e.ErrorAuthCheckTokenTimeout
			} else if claims.Authority == 0 {
				code = e.ErrorAuthInsufficientAuthority
			} else {
				c.Set("user_id", claims.ID)
				c.Set("username", claims.Username)
			}
		}

		if code != e.SUCCESS {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status": code,
				"msg":    e.GetMsg(code),
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

func extractToken(c *gin.Context) string {
	auth := c.GetHeader("Authorization")
	if strings.HasPrefix(auth, "Bearer ") {
		return strings.TrimPrefix(auth, "Bearer ")
	}
	if token := c.GetHeader("Cookie"); token != "" {
		return token
	}
	return ""
}
