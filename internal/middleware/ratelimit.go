package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type redisSlidingWindow struct {
	client   *redis.Client
	key      string
	window   time.Duration
	maxCount int64
}

func NewRedisSlidingWindow(client *redis.Client, key string, window time.Duration, maxCount int64) *redisSlidingWindow {
	return &redisSlidingWindow{
		client:   client,
		key:      key,
		window:   window,
		maxCount: maxCount,
	}
}

func (r *redisSlidingWindow) Allow(ctx context.Context) (bool, error) {
	now := time.Now().UnixMilli()
	windowStart := now - r.window.Milliseconds()

	pipe := r.client.Pipeline()
	pipe.ZRemRangeByScore(ctx, r.key, "0", fmt.Sprintf("%d", windowStart))
	pipe.ZAdd(ctx, r.key, redis.Z{Score: float64(now), Member: now})
	pipe.Expire(ctx, r.key, r.window)
	_, err := pipe.Exec(ctx)
	if err != nil {
		return false, err
	}

	total, err := r.client.ZCard(ctx, r.key).Result()
	if err != nil {
		return false, err
	}

	return total <= r.maxCount, nil
}

func RateLimit(client *redis.Client, maxReq int, windowSec int, keyPrefix string) gin.HandlerFunc {
	limiter := NewRedisSlidingWindow(client, keyPrefix, time.Duration(windowSec)*time.Second, int64(maxReq))

	return func(c *gin.Context) {
		allowed, err := limiter.Allow(c.Request.Context())
		if err != nil {
			c.Next()
			return
		}

		if !allowed {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"status": 429,
				"msg":    "too many requests",
			})
			return
		}

		c.Next()
	}
}