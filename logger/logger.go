package logger

import (
	"BlockMe/config"
	"fmt"
	"io"
	"log/slog"
	"os"
)

func Init() error {
	hostname, err := os.Hostname()
	if err != nil {
		fmt.Printf("Unable to get hostname: %v\n", err)
		hostname = "unknown"
	}

	var logWriter io.Writer

	if config.Env.LogLocation == "file" {
		logWriter, err = os.OpenFile("blockme.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			return err
		}
	} else {
		logWriter = os.Stdout
	}

	var logHandler slog.Handler
	switch config.Env.LogFormat {
	case "json":
		logHandler = slog.NewJSONHandler(logWriter, nil)
	case "text":
		logHandler = slog.NewTextHandler(logWriter, nil)
	case "text_pretty":
		logHandler = slog.NewTextHandler(logWriter, nil)
	default:
		logHandler = slog.NewTextHandler(logWriter, nil)
	}

	logger := slog.New(logHandler)
	logger.With([]slog.Attr{slog.String("hostname", hostname)})

	slog.SetDefault(logger)
	switch config.Env.LogLevel {
	case "DEBUG":

		slog.SetLogLoggerLevel(slog.LevelDebug)
	case "INFO":
		slog.SetLogLoggerLevel(slog.LevelInfo)
	case "WARN":
		slog.SetLogLoggerLevel(slog.LevelWarn)
	case "ERROR":
		slog.SetLogLoggerLevel(slog.LevelError)
	default:
		slog.SetLogLoggerLevel(slog.LevelInfo)
	}

	return nil
}
