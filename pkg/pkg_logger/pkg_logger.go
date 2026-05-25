package pkg_logger

import (
	"log/slog"
	"os"
	"path"
	"time"
)

var Logger *slog.Logger

func init() {
	if Logger != nil {
		return
	}
	logFilePath := ""
	if dir, err := os.Getwd(); err == nil {
		logFilePath = dir + "/logs/"
	}

	if err := os.MkdirAll(logFilePath, 0755); err != nil {
		panic(err)
	}

	logFileName := time.Now().Format("2006-01-02") + ".log"
	fileName := path.Join(logFilePath, logFileName)

	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}

	handler := slog.NewJSONHandler(file, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	Logger = slog.New(handler)
}