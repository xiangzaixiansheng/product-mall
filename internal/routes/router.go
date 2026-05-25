package routes

import (
	"net/http"
	"time"

	api "product-mall/internal/api/v1"
	"product-mall/internal/middleware"
	"product-mall/pkg/db"

	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())

	r.Use(gzip.Gzip(gzip.DefaultCompression))
	r.Use(middleware.SecurityHeaders())
	r.Use(middleware.WithRequsetId())
	r.Use(middleware.Timeout(30 * time.Second))
	r.Use(middleware.Logger())

	store := cookie.NewStore([]byte("something-very-secret"))
	r.Use(sessions.Sessions("mysession", store))

	r.GET("/health", api.HealthCheck)

	redisClient := db.GetRedisClient()
	if redisClient != nil {
		r.Use(middleware.RateLimit(redisClient, 100, 60, "ratelimit:"))
	}

	v1 := r.Group("api/v1")
	{
		v1.GET("ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, "success")
		})

		v1.POST("user/register", api.UserRegister)
		v1.POST("user/login", api.UserLogin)

		authed := v1.Group("/")
		authed.Use(middleware.JWT())
		{
			authed.PUT("user", api.UserUpdate)
			authed.POST("user/sending-email", api.SendEmail)
			authed.POST("product", api.CreateProduct)
			authed.POST("carts/:id", api.CreateCart)
			authed.GET("carts/:id", api.ShowList)
			authed.PUT("carts/:id", api.UpdateCart)
			authed.DELETE("carts/:id", api.DeleteCart)
			authed.POST("addresses", api.CreateAddress)
			authed.GET("addresses/:id", api.ShowAddresses)
			authed.PUT("addresses/:id", api.UpdateAddress)
			authed.DELETE("addresses/:id", api.DeleteAddress)
		}
	}
	return r
}