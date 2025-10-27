package logger

import (
	"log/slog"
	"os"

	"gopkg.in/natefinch/lumberjack.v2"
)

var Logger *slog.Logger

func Init(level string) {
	var slogLevel slog.Level

	switch level {
	case "debug":
		slogLevel = slog.LevelDebug
	case "info":
		slogLevel = slog.LevelInfo
	case "warn":
		slogLevel = slog.LevelWarn
	case "error":
		slogLevel = slog.LevelError
	default:
		slogLevel = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level: slogLevel,
	}

	logWriter := &lumberjack.Logger{
		Filename:   "logs/app.log",
		MaxSize:    100,
		MaxBackups: 3,
		MaxAge:     28,
		Compress:   true,
	}

	handler := slog.NewJSONHandler(logWriter, opts)
	Logger = slog.New(handler)
	
	os.MkdirAll("logs", 0755)
}

func Info(msg string, fields ...any) {
	if Logger != nil {
		Logger.Info(msg, fields...)
	}
}

func Error(msg string, fields ...any) {
	if Logger != nil {
		Logger.Error(msg, fields...)
	}
}

func Debug(msg string, fields ...any) {
	if Logger != nil {
		Logger.Debug(msg, fields...)
	}
}

func Warn(msg string, fields ...any) {
	if Logger != nil {
		Logger.Warn(msg, fields...)
	}
}