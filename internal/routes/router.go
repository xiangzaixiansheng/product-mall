package routes

import (
	"net/http"
	"os"
	"time"

	api "product-mall/internal/api/v1"
	"product-mall/internal/middleware"
	"product-mall/pkg/db"

	_ "product-mall/docs/swagger"

	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())

	r.Use(gzip.Gzip(gzip.DefaultCompression))
	r.Use(middleware.SecurityHeaders())
	r.Use(middleware.WithRequsetId())
	r.Use(middleware.Timeout(30 * time.Second))
	r.Use(middleware.Logger())

	sessionSecret := os.Getenv("SESSION_SECRET")
	if sessionSecret == "" {
		sessionSecret = "dev-session-secret-change-me"
	}
	store := cookie.NewStore([]byte(sessionSecret))
	r.Use(sessions.Sessions("mysession", store))

	r.GET("/health", api.HealthCheck)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	redisClient := db.GetRedisClient()
	if redisClient != nil {
		r.Use(middleware.RateLimit(redisClient, 100, 60, "ratelimit:"))
	}

	v1 := r.Group("api/v1")
	{
		v1.GET("ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, "success")
		})

		// 公开接口
		v1.POST("user/register", api.UserRegister)
		v1.POST("user/login", api.UserLogin)
		v1.GET("categories", api.ListCategories)
		v1.GET("products", api.SearchProducts)

		// 需要认证的接口
		authed := v1.Group("/")
		authed.Use(middleware.JWT())
		{
			// 用户
			authed.PUT("user", api.UserUpdate)
			authed.POST("user/sending-email", api.SendEmail)

			// 分类管理
			authed.POST("categories", api.CreateCategory)
			authed.PUT("categories/:id", api.UpdateCategory)
			authed.DELETE("categories/:id", api.DeleteCategory)

			// 商品
			authed.POST("product", api.CreateProduct)

			// 购物车
			authed.POST("carts/:id", api.CreateCart)
			authed.GET("carts/:id", api.ShowList)
			authed.PUT("carts/:id", api.UpdateCart)
			authed.DELETE("carts/:id", api.DeleteCart)

			// 地址
			authed.POST("addresses", api.CreateAddress)
			authed.GET("addresses/:id", api.ShowAddresses)
			authed.PUT("addresses/:id", api.UpdateAddress)
			authed.DELETE("addresses/:id", api.DeleteAddress)

			// 订单
			authed.POST("orders", api.CreateOrder)
			authed.GET("orders", api.ListOrders)
			authed.GET("orders/:id", api.GetOrder)
			authed.PUT("orders/:id/pay", api.PayOrder)
			authed.PUT("orders/:id/cancel", api.CancelOrder)
		}
	}
	return r
}
