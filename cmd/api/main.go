package main

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"url-shortener/internal/cache"
	appcfg "url-shortener/internal/config"
	"url-shortener/internal/logger"
	mongorepo "url-shortener/internal/repo/mongo"
	"url-shortener/internal/server"
)

func main() {

	// config init
	var cfg = appcfg.Get()

	// Logger'ı başlat
	loggerInstance := logger.GetLogger()
	defer func() {
		if l, ok := logger.GetLogger().(*logger.loggerImpl); ok {
			_ = l.zapLogger.Sync() // Logları temizler
		}
	}()

	// Infra compose (DB’ler)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mcli, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		log.Fatal("mongo connect:", err)
	}
	if err := mcli.Ping(ctx, nil); err != nil {
		log.Fatal("mongo ping:", err)
	}

	db := mcli.Database(cfg.MongoDB)
	urlRepo := mongorepo.NewURLRepo(db)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	redis := cache.NewRedis(cfg.RedisAddr, cfg.RedisPassword, cfg.RedisDB)
	if err := redis.Rdb.Ping(ctx).Err(); err != nil {
		log.Fatal(err)
	}

	loggerInstance.Info("starting server")
	srv := server.New(server.Options{Port: cfg.Port, BaseURL: cfg.BaseURL, Repo: urlRepo, Cache: redis, Logger: loggerInstance})

	if err := srv.Start(ctx); err != nil {
		loggerInstance.Error("server stopped with error:", zap.Error(err))
	}

	_ = mcli.Disconnect(context.Background())
}
