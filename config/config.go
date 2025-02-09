package config

import (
	"os"
)

type Config struct {
	EventDBURL       string
	ProjectionDBURL  string
	KafkaURL         string
	APIPort          string
	EventAPIPort     string
	DatabaseName     string
	DatabaseUser     string
	DatabasePassword string
	DatabaseHost     string
	DatabasePort     string
}

func Load() Config {
	return Config{
		EventDBURL:       os.Getenv("EVENT_DB_URL"),
		ProjectionDBURL:  os.Getenv("PROJECTION_DB_URL"),
		KafkaURL:         os.Getenv("KAFKA_URL"),
		EventAPIPort:     os.Getenv("EVENT_API_PORT"),
		APIPort:          os.Getenv("API_PORT"),
		DatabaseName:     os.Getenv("DATABASE_NAME"),
		DatabaseUser:     os.Getenv("DATABASE_USER"),
		DatabasePassword: os.Getenv("DATABASE_PASSWORD"),
		DatabaseHost:     os.Getenv("DATABASE_HOST"),
		DatabasePort:     os.Getenv("DATABASE_PORT"),
	}
}
