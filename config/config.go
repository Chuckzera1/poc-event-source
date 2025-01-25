package config

import (
	"os"
)

type Config struct {
	EventDBURL      string
	ProjectionDBURL string
	KafkaURL        string
}

func Load() Config {
	return Config{
		EventDBURL:      os.Getenv("EVENT_DB_URL"),
		ProjectionDBURL: os.Getenv("PROJECTION_DB_URL"),
		KafkaURL:        os.Getenv("KAFKA_URL"),
	}
}
