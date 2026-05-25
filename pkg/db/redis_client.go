package db

import (
	"context"
	"product-mall/conf"

	"github.com/redis/go-redis/v9"
	"github.com/go-redis/redismock/v9"
)

var (
	client *redis.Client
	Mocker redismock.ClientMock
)

func InitMockClient() {
	cli, mock := redismock.NewClientMock()
	client = cli
	Mocker = mock
}

func InitRedis(ctx context.Context) (err error) {
	client = redis.NewClient(&redis.Options{
		Addr:     conf.RedisAddr,
		Password: conf.RedisPw,
		DB:       0,
	})
	_, err = client.Ping(ctx).Result()
	if err != nil {
		return
	}
	return err
}

func GetRedisClient() *redis.Client {
	return client
}
