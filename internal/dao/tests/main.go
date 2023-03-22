//go:build database
// +build database

// To override application configuration for integration tests, create config/test.env file.

package tests

import (
	"context"
	"os"
	"testing"

	_ "consoledot-go-template/internal/dao/pgx"
)

// truncate and seed database tables
func reset() {
	DbSeed()
}

func TestMain(t *testing.M) {
	ctx := context.Background()
	InitEnvironment(ctx, "../../../config/test.env")
	defer CloseEnvironment(ctx)
	defer DbDrop()

	DbDrop()
	DbMigrate()
	reset()
	exitVal := t.Run()
	os.Exit(exitVal)
}
