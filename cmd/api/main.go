package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"poc-event-source/config"
	"poc-event-source/internal/api"
	"poc-event-source/internal/api/routes"
	"poc-event-source/internal/infrastructure"
	"time"
)

func init() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("erro ao carregar o arquivo .env: %v", err)
	}
}

func main() {
	cfg := config.Load()

	fmt.Println("Pass -> ", cfg.DatabasePassword)
	_, err := infrastructure.NewGormDB(fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DatabaseHost,
		cfg.DatabasePort,
		cfg.DatabaseUser,
		cfg.DatabasePassword,
		cfg.DatabaseName,
	))
	if err != nil {
		log.Fatalf("Error connecting to events db: %v", err)
	}

	_, err = infrastructure.NewGormDB(cfg.ProjectionDBURL)
	if err != nil {
		log.Printf("Error connecting to projection db: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 168*time.Hour)
	_, err = infrastructure.Nats(cfg.BrokerURL, cfg.BrokerStreamName, cfg.BrokerSubjects, ctx, cancel)
	if err != nil {
		log.Fatalf("Error connecting to events broker: %v", err)
	}

	err = api.StartAPI(cfg, routes.SetupUserRouter)
	if err != nil {
		log.Fatalf("Error starting api: %v", err)
	}
}
