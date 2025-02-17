package testutils

import (
	"context"
	"fmt"
	"github.com/testcontainers/testcontainers-go"
	"time"

	"github.com/docker/go-connections/nat"
	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go/wait"

	"gorm.io/gorm"
)

type TestDatabase struct {
	Host      string
	Port      string
	User      string
	Password  string
	Database  string
	container testcontainers.Container
	GormDB    *gorm.DB
}

func NewTestDatabase(ctx context.Context, model ...interface{}) (*TestDatabase, error) {
	container, err := startPostgresContainer(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to start PostgreSQL container: %w", err)
	}

	dbConfig, err := getContainerConfig(ctx, container)
	if err != nil {
		err := container.Terminate(ctx)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("failed to get container config: %w", err)
	}

	gormDB, err := connectWithGorm(dbConfig)
	if err != nil {
		err := container.Terminate(ctx)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("failed to connect with GORM: %w", err)
	}

	err = gormDB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";").Error
	if err != nil {
		return nil, fmt.Errorf("failed to create uuid-ossp extension: %v", err)
	}

	if len(model) > 0 {
		if err := autoMigratemodel(gormDB, model); err != nil {
			err := container.Terminate(ctx)
			if err != nil {
				return nil, err
			}
			return nil, fmt.Errorf("failed to run AutoMigrate: %w", err)
		}
	}

	return &TestDatabase{
		Host:      dbConfig.Host,
		Port:      dbConfig.Port,
		User:      dbConfig.User,
		Password:  dbConfig.Password,
		Database:  dbConfig.Database,
		container: container,
		GormDB:    gormDB,
	}, nil
}

func startPostgresContainer(ctx context.Context) (testcontainers.Container, error) {
	req := testcontainers.ContainerRequest{
		Image:        "postgres:16-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "testuser",
			"POSTGRES_PASSWORD": "testpass",
			"POSTGRES_DB":       "testdb",
		},
		WaitingFor: wait.ForSQL("5432/tcp", "postgres", func(host string, port nat.Port) string {
			return fmt.Sprintf("host=%s port=%s user=testuser password=testpass dbname=testdb sslmode=disable", host, port.Port())
		}).WithStartupTimeout(30 * time.Second),
	}

	return testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
}

func autoMigratemodel(db *gorm.DB, model []interface{}) error {
	return db.AutoMigrate(model...)
}

func (db *TestDatabase) Cleanup(ctx context.Context) error {
	return db.container.Terminate(ctx)
}
