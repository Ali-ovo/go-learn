package main

import (
	"time"

	"go.uber.org/zap"
)

func NewLogger() (*zap.Logger, error) {
	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{
		"./myproject.log",
		"stdout",
	}

	return cfg.Build()
}

func main() {
	// logger, _ := zap.NewProduction()
	logger, err := NewLogger()
	if err != nil {
		panic(err)
	}

	defer logger.Sync() // flushes buffer, if any

	url := "http://example.com"

	// logger.Info("failed to fetch URL",
	// 	zap.String("url", url),
	// 	zap.Int("num", 3),
	// )

	sugar := logger.Sugar()
	sugar.Infow("failed to fetch URL",
		// Structured context as loosely typed key-value pairs.
		"url", url,
		"attempt", 3,
		"backoff", time.Second,
	)
}
