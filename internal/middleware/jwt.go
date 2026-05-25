package middleware

import (
	"net/http"
	"product-mall/conf"
	util "product-mall/internal/tools"
	"product-mall/pkg/e"
	"time"

	"github.com/gin-gonic/gin"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var code int
		var data any
		code = 200
		token := c.GetHeader("Cookie")

		if conf.ENV == "dev" {
			//测试环境走下去
			c.Next()
			return
		}

		if token == "" {
			code = 404
		} else {
			claims, err := util.ParseToken(token)
			if err != nil {
				code = e.ErrorAuthCheckTokenFail
			} else if claims.ExpiresAt.Before(time.Now()) {
				code = e.ErrorAuthCheckTokenTimeout
			}
		}
		if code != e.SUCCESS {
			c.JSON(http.StatusOK, gin.H{
				"status": code,
				"msg":    e.GetMsg(code),
				"data":   data,
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

//JWTAdmin token验证中间件
func JWTAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		var code int
		var data any
		token := c.GetHeader("Cookie")
		if conf.ENV == "dev" {
			//测试环境走下去
			c.Next()
			return
		}
		if token == "" {
			code = e.InvalidParams
		} else {
			claims, err := util.ParseToken(token)
			if err != nil {
				code = e.ErrorAuthCheckTokenFail
			} else if claims.ExpiresAt.Before(time.Now()) {
				code = e.ErrorAuthCheckTokenTimeout
			} else if claims.Authority == 0 {
				code = e.ErrorAuthInsufficientAuthority
			}
		}
		if code != e.SUCCESS {
			c.JSON(http.StatusOK, gin.H{
				"status": code,
				"msg":    e.GetMsg(code),
				"data":   data,
			})
			c.Abort()
			return
		}
		//走下去
		c.Next()
	}
}
