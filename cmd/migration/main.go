package main

import (
	"context"
	appcfg "github.com/emrealsandev/Url-Shortener/internal/config"
	"github.com/emrealsandev/Url-Shortener/internal/migration"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

func main() {
	_ = godotenv.Load()

	var cfg appcfg.Config
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatal("config: ", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		log.Fatal("mongo connect: ", err)
	}
	defer client.Disconnect(context.Background())

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal("mongo ping: ", err)
	}

	db := client.Database(cfg.MongoDB)
	m := migration.New(db)
	if err := m.RunAll(ctx); err != nil {
		log.Fatal("migration failed: ", err)
	}

	log.Println("migration OK")
}
