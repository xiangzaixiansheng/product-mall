// @title Product Mall API
// @version 1.0
// @description 电商商城后端 API 服务
// @host localhost:3000
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"product-mall/conf"
	"product-mall/internal/routes"
	"product-mall/pkg/db"
)

func main() {
	conf.Init()
	db.Database(conf.MysqlpathRead, conf.MysqlpathWrite)

	r := routes.NewRouter()
	srv := &http.Server{
		Addr:         conf.HttpPort,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		slog.Info("server starting", "addr", conf.HttpPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("listen failed", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("server forced to shutdown", "error", err)
	}

	if sqlDB, err := db.GetDB().DB(); err == nil {
		if err := sqlDB.Close(); err != nil {
			slog.Error("mysql close failed", "error", err)
		}
	}

	if redisClient := db.GetRedisClient(); redisClient != nil {
		if err := redisClient.Close(); err != nil {
			slog.Error("redis close failed", "error", err)
		}
	}

	slog.Info("server exited")
}
