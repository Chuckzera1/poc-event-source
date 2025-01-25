package api

import (
	"log"
	"poc-event-source/config"
	"poc-event-source/internal/infrastructure"
)

func main() {
	cfg := config.Load()

	_, err := infrastructure.NewGormDB(cfg.EventDBURL)
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco de eventos: %v", err)
	}

	_, err = infrastructure.NewGormDB(cfg.ProjectionDBURL)
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco de projeções: %v", err)
	}
}
