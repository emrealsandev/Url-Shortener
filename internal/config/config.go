package config

type Config struct {
	Port    string `envconfig:"PORT" default:"8080"`
	BaseURL string `envconfig:"BASE_URL" required:"true"`

	MongoURI string `envconfig:"MONGO_URI" required:"true"`
	MongoDB  string `envconfig:"MONGO_DB" default:"shortener"`
}
