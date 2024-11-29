// internal/logger/logger.go

package logger

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewLogger initializes the logger
func NewLogger(environment string) (*zap.Logger, error) {
	var logger *zap.Logger
	var err error

	switch environment {
	case "production":
		logger, err = productionLogger()
	default:
		logger, err = developmentLogger()
	}

	if err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %v", err)
	}

	return logger, nil
}

// developmentLogger creates a development logger
func developmentLogger() (*zap.Logger, error) {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	return config.Build()
}

// productionLogger creates a production logger
func productionLogger() (*zap.Logger, error) {
	config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	return config.Build()
}
