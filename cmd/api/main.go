package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	appcfg "url-shortener/internal/config"
	mongorepo "url-shortener/internal/repo/mongo"
	"url-shortener/internal/server"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	_ = godotenv.Load()

	var cfg appcfg.Config
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatal(err)
	}

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
	urls := db.Collection("urls")
	if err := mongorepo.EnsureIndexes(ctx, urls); err != nil {
		log.Fatal("indexes:", err)
	}
	urlRepo := mongorepo.NewURLRepo(urls)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	srv := server.New(server.Options{Port: cfg.Port, BaseURL: cfg.BaseURL, Repo: urlRepo})

	if err := srv.Start(ctx); err != nil {
		log.Fatal("server stopped with error:", err)
	}

	_ = mcli.Disconnect(context.Background())
}
