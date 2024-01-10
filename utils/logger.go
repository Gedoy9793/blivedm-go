package utils

import (
	"context"
	"github.com/sirupsen/logrus"
)

func GetLoggerFromContext(ctx context.Context) logrus.FieldLogger {
	logger, ok := ctx.Value("Logger").(logrus.FieldLogger)
	if !ok {
		logger = logrus.StandardLogger()
	}
	return logger
}

func SetLoggerToContext(ctx context.Context, logger logrus.FieldLogger) context.Context {
	return context.WithValue(ctx, "Logger", logger)
}
