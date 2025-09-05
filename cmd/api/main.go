package main

import (
	"context"
	"go.uber.org/zap"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"url-shortener/internal/logger"

	appcfg "url-shortener/internal/config"
	mongorepo "url-shortener/internal/repo/mongo"
	"url-shortener/internal/server"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {

	// config init
	var cfg = appcfg.Get()

	// logger init
	defer logger.L().Sync()

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

	logger.L().Info("starting server")
	srv := server.New(server.Options{Port: cfg.Port, BaseURL: cfg.BaseURL, Repo: urlRepo})

	if err := srv.Start(ctx); err != nil {
		logger.L().Error("server stopped with error:", zap.Error(err))
	}

	_ = mcli.Disconnect(context.Background())
}
