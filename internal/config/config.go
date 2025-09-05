package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"sync"
)

const ENVIRONMENT_LOCAL = "local"
const ENVIRONMENT_PROD = "prod"

type Config struct {
	Port    string `envconfig:"PORT" default:"8080"`
	BaseURL string `envconfig:"BASE_URL" required:"true"`

	MongoURI     string `envconfig:"MONGO_URI" required:"true"`
	MongoDB      string `envconfig:"MONGO_DB" default:"shortener"`
	Environment  string `envconfig:"APP_ENVIRONMENT" default:"dev"`
	SequenceSalt string `envconfig:"SEQUENCE_SALT" default:"_"`
}

var (
	cfg  *Config
	once sync.Once
)

func Load() *Config {
	once.Do(func() {
		if err := godotenv.Load(); err != nil {
			log.Println("⚠️ .env yüklenemedi:", err)
		}
		cfg = &Config{
			Port:         os.Getenv("PORT"),
			BaseURL:      os.Getenv("BASE_URL"),
			MongoURI:     os.Getenv("MONGO_URI"),
			MongoDB:      os.Getenv("MONGO_DB"),
			Environment:  os.Getenv("APP_ENVIRONMENT"),
			SequenceSalt: os.Getenv("SEQUENCE_SALT"),
		}
	})
	return cfg
}

func Get() *Config {
	return Load()
}
