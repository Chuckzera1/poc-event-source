package config

import (
	"github.com/nats-io/nats.go"
	"os"
	"strings"
)

type Config struct {
	EventDBURL       string
	ProjectionDBURL  string
	BrokerURL        string
	BrokerStreamName string
	BrokerSubjects   []string
	APIPort          string
	EventAPIPort     string
	DatabaseName     string
	DatabaseUser     string
	DatabasePassword string
	DatabaseHost     string
	DatabasePort     string
}

func Load() Config {
	brokerUrl := os.Getenv("BROKER_URL")
	if brokerUrl == "" {
		brokerUrl = nats.DefaultURL
	}
	subjects := os.Getenv("BROKER_SUBJECTS")
	if subjects == "" {
		subjects = "*"
	}
	brokerSubjects := strings.Split(subjects, ",")
	return Config{
		EventDBURL:       os.Getenv("EVENT_DB_URL"),
		ProjectionDBURL:  os.Getenv("PROJECTION_DB_URL"),
		BrokerURL:        brokerUrl,
		BrokerSubjects:   brokerSubjects,
		BrokerStreamName: os.Getenv("BROKER_STREAM_NAME"),
		EventAPIPort:     os.Getenv("EVENT_API_PORT"),
		APIPort:          os.Getenv("API_PORT"),
		DatabaseName:     os.Getenv("DATABASE_NAME"),
		DatabaseUser:     os.Getenv("DATABASE_USER"),
		DatabasePassword: os.Getenv("DATABASE_PASSWORD"),
		DatabaseHost:     os.Getenv("DATABASE_HOST"),
		DatabasePort:     os.Getenv("DATABASE_PORT"),
	}
}
