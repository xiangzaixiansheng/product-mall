package cache

import (
	"context"
	"strconv"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

type MyRedis struct {
	Client *redis.Client
}

var once sync.Once

var RedisClient *MyRedis = new(MyRedis)

func GetInstance() *MyRedis {
	return RedisClient
}

func NewRedis(RedisAddr string, RedisDbName string, RedisPw string) *redis.Client {
	db, _ := strconv.ParseUint(RedisDbName, 10, 64)
	myRedis := redis.NewClient(&redis.Options{
		Addr: RedisAddr,
		DB:   int(db),
	})
	ctx := context.Background()
	_, err := myRedis.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}

	once.Do(func() {
		RedisClient.Client = myRedis
	})

	return myRedis
}

func (mr *MyRedis) Set(key string, value any, ttl time.Duration) {
	mr.Client.Set(context.Background(), key, value, ttl)
}

func (mr MyRedis) Get(key string) *redis.StringCmd {
	return mr.Client.Get(context.Background(), key)
}

func (mr MyRedis) Incr(key string) *redis.IntCmd {
	return mr.Client.Incr(context.Background(), key)
}

func (mr MyRedis) ZIncrBy(key string, increment float64, member string) *redis.FloatCmd {
	return mr.Client.ZIncrBy(context.Background(), key, increment, member)
}

func (mr *MyRedis) Lock(key string, ttl time.Duration) (bool, error) {
	return mr.Client.SetNX(context.Background(), key, 1, ttl).Result()
}

func (mr *MyRedis) Unlock(key string) error {
	return mr.Client.Del(context.Background(), key).Err()
}
