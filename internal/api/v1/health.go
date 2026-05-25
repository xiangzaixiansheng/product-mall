package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"product-mall/pkg/db"
	"product-mall/cache"
)

type HealthResponse struct {
	Status   string            `json:"status"`
	Services map[string]string `json:"services"`
}

func HealthCheck(c *gin.Context) {
	response := HealthResponse{
		Status:   "ok",
		Services: make(map[string]string),
	}

	// Check MySQL
	sqlDB, err := db.DB.DB()
	if err != nil || sqlDB.Ping() != nil {
		response.Services["mysql"] = "down"
		response.Status = "degraded"
	} else {
		response.Services["mysql"] = "ok"
	}

	// Check Redis
	redisClient := cache.GetInstance().Client
	if redisClient == nil {
		response.Services["redis"] = "down"
		response.Status = "degraded"
	} else if _, err := redisClient.Ping(c.Request.Context()).Result(); err != nil {
		response.Services["redis"] = "down"
		response.Status = "degraded"
	} else {
		response.Services["redis"] = "ok"
	}

	statusCode := http.StatusOK
	if response.Status != "ok" {
		statusCode = http.StatusServiceUnavailable
	}

	c.JSON(statusCode, response)
}