package logger

import (
	"go.uber.org/zap"
	"os"
	"url-shortener/internal/config"
)

var log *zap.Logger

func Init() {
	var err error
	switch os.Getenv("APP_ENVIRONMENT") {
	case config.ENVIRONMENT_PROD:
		log, err = zap.NewProduction()
	case config.ENVIRONMENT_LOCAL:
		log, err = zap.NewDevelopment()
	}

	if err != nil {
		panic(err)
	}
}

func L() *zap.Logger {
	if log == nil {
		Init()
	}
	return log
}
