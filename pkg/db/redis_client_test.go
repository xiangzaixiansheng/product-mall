package db

import (
	"context"
	"errors"
	"product-mall/conf"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/redis/go-redis/v9"
	. "github.com/smartystreets/goconvey/convey"
)

func init() {
	conf.Init()
	InitMockClient()
}

func TestInitRedis(t *testing.T) {
	Convey("Test InitRedis", t, func() {
		Convey("should success", func() {
			patches := gomonkey.NewPatches()
			defer patches.Reset()
			Mocker.ExpectPing().SetVal("Pong")
			patches.ApplyFuncReturn(redis.NewClient, client)
			err := InitRedis(context.Background())
			So(err, ShouldBeNil)
		})
		Convey("should fail", func() {
			patches := gomonkey.NewPatches()
			defer patches.Reset()
			Mocker.ExpectPing().SetErr(errors.New("test error"))
			patches.ApplyFuncReturn(redis.NewClient, client)
			err := InitRedis(context.Background())
			So(err, ShouldNotBeNil)
		})
	})
}
