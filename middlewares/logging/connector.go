package logging

import (
	"go.uber.org/zap"
)

// NewLogger создаёт и возвращает новый zap.Logger
func NewLogger() *zap.Logger {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	return logger
}
