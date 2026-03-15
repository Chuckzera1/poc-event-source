package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"poc-event-source/config"
	"poc-event-source/internal/api"
	apimessaging "poc-event-source/internal/api/messaging"
	"poc-event-source/internal/api/routes"
	usecaseevent "poc-event-source/internal/application/usecase/event"
	usecaseuser "poc-event-source/internal/application/usecase/user"
	"poc-event-source/internal/application/utils"
	"poc-event-source/internal/infrastructure"
	eventrepo "poc-event-source/internal/repository/event"
	userrepo "poc-event-source/internal/repository/user"
)

func init() {
	if err := godotenv.Load(".env"); err != nil {
		log.Printf("warning: .env file not found: %v", err)
	}
}

func main() {
	cfg := config.Load()

	eventsDB, err := infrastructure.NewGormDB(fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DatabaseHost,
		cfg.DatabasePort,
		cfg.DatabaseUser,
		cfg.DatabasePassword,
		cfg.DatabaseName,
	))
	if err != nil {
		log.Fatalf("Error connecting to events db: %v", err)
	}

	projectionDB, err := infrastructure.NewGormDB(cfg.ProjectionDBURL)
	if err != nil {
		log.Fatalf("Error connecting to projection db: %v", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	broker, err := infrastructure.Nats(cfg.BrokerURL, cfg.BrokerStreamName, cfg.BrokerSubjects, ctx, cancel)
	if err != nil {
		log.Fatalf("Error connecting to events broker: %v", err)
	}

	eventRepo := eventrepo.NewEventRepository(eventsDB)
	userRepo := userrepo.NewUserRepository(projectionDB)
	pwdUtil := utils.NewPasswordUtils(bcrypt.DefaultCost)

	eventHandler := usecaseevent.NewMainHandler(eventRepo, broker)
	createUserUC := usecaseuser.NewCreateUserUseCase(eventHandler)

	userBroker := apimessaging.NewUserBroker(broker, userRepo, pwdUtil)
	if err := userBroker.Subscribe(); err != nil {
		log.Fatalf("Error starting user subscriber: %v", err)
	}

	userHandler := routes.NewUserHandler(createUserUC)
	if err := api.StartAPI(cfg, userHandler.SetupUserRouter); err != nil {
		log.Fatalf("Error starting api: %v", err)
	}
}
