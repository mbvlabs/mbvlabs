package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"mbvlabs/config"
	"mbvlabs/migrations"

	"github.com/mbvlabs/andurel/pkg/env"
	"github.com/mbvlabs/andurel/pkg/storage"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	if err := env.Load(); err != nil {
		return err
	}

	ctx := context.Background()

	db, err := openMigrationDatabase(ctx)
	if err != nil {
		return err
	}
	defer db.Close()

	fmt.Println("Applying migrations...")
	if err := storage.RunMigrations(ctx, db, migrations.Migrations, "."); err != nil {
		return err
	}

	fmt.Println("Migrations complete!")
	return nil
}

func openMigrationDatabase(ctx context.Context) (*storage.Postgres, error) {
	if migrationURL := strings.TrimSpace(os.Getenv("DATABASE_MIGRATION_URL")); migrationURL != "" {
		db, err := storage.NewPostgres(ctx, storage.DefaultConfig(), storage.WithDatabaseURL(migrationURL))
		if err != nil {
			// A connection/parse error can contain the credential-bearing URL.
			return nil, errors.New("failed to connect using DATABASE_MIGRATION_URL; check the PostgreSQL URL and connectivity")
		}
		return db, nil
	}

	if os.Getenv("RELEASE_COMMAND") == "1" {
		return nil, errors.New("DATABASE_MIGRATION_URL is required for the Fly release command")
	}

	cfg, err := config.NewDatabase()
	if err != nil {
		return nil, err
	}
	db, err := storage.NewPostgres(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	return db, nil
}
