package logger

import (
	"fmt"
	"log/slog"
	"os"
)

func Init() error {
	hostname, err := os.Hostname()
	if err != nil {
		fmt.Printf("Unable to get hostname: %v\n", err)
		hostname = "unknown"
	}

	accessLog, err := os.OpenFile("access.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}

	logger := slog.New(slog.NewJSONHandler(accessLog, nil))
	logger.With([]slog.Attr{slog.String("hostname", hostname)})

	slog.SetDefault(logger)
	slog.SetLogLoggerLevel(slog.LevelInfo)

	return nil
}
