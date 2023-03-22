//go:build database
// +build database

package tests

import (
	"consoledot-go-template/internal/config"
	"consoledot-go-template/internal/db"
	"consoledot-go-template/internal/db/seeds"
	"consoledot-go-template/internal/logging"
	"context"
	"fmt"
)

func InitEnvironment(ctx context.Context, envPath string) {
	config.Initialize("config/test.env", envPath)
	logger, _ := logging.InitializeLogger()
	ctx = logging.WithLogger(ctx, &logger)

	err := db.Initialize(context.Background(), "integration")
	if err != nil {
		panic(fmt.Errorf("cannot connect to database: %w (integration schema)", err))
	}
}

func CloseEnvironment(ctx context.Context) {
	db.Close()
}

func DbDrop() {
	err := seeds.Seed(context.Background(), "drop_integration")
	if err != nil {
		panic(err)
	}
}

func DbMigrate() {
	err := db.Migrate(context.Background(), "integration")
	if err != nil {
		panic(err)
	}
}

func DbSeed() {
	err := seeds.Seed(context.Background(), "integration")
	if err != nil {
		panic(err)
	}
}
