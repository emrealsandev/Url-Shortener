package logger

import (
	"github.com/emrealsandev/Url-Shortener/internal/config"
	"os"

	"go.uber.org/zap"
)

type loggerImpl struct {
	zapLogger *zap.Logger
}

var singleton Logger

func newLogger() {
	var zapLogger *zap.Logger
	var err error

	switch os.Getenv("APP_ENVIRONMENT") {
	case config.ENVIRONMENT_PROD:
		zapLogger, err = zap.NewProduction()
	case config.ENVIRONMENT_LOCAL:
		zapLogger, err = zap.NewDevelopment()
	default:
		zapLogger, err = zap.NewDevelopment()
	}

	if err != nil {
		panic(err)
	}

	singleton = &loggerImpl{zapLogger: zapLogger}
}

func GetLogger() Logger {
	if singleton == nil {
		newLogger()
	}
	return singleton
}

func (l *loggerImpl) Debug(msg string, fields ...any) {
	l.zapLogger.Sugar().Debugw(msg, fields...)
}

func (l *loggerImpl) Info(msg string, fields ...any) {
	l.zapLogger.Sugar().Infow(msg, fields...)
}

func (l *loggerImpl) Warn(msg string, fields ...any) {
	l.zapLogger.Sugar().Warnw(msg, fields...)
}

func (l *loggerImpl) Error(msg string, fields ...any) {
	l.zapLogger.Sugar().Errorw(msg, fields...)
}

func (l *loggerImpl) Sync() {
	l.zapLogger.Sync()
}
