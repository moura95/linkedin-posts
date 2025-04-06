package db

import (
	"database/sql"
	"embed"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

func RunMigrations(db *sql.DB) error {
	log.Println("Initializing migrations...")

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("error connecting to database: %w", err)
	}

	source, err := iofs.New(migrationsFS, "migrations")
	if err != nil {
		return fmt.Errorf("Erorr to load migrations: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", source, "postgres", driver)
	if err != nil {
		return fmt.Errorf("Error to instance migration: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("Error to execute migrations: %w", err)
	}

	log.Println("Migrations executed successfully")
	return nil
}

func RunMigrationsWithRetry(db *sql.DB, maxRetries int) error {
	var err error
	for i := 0; i < maxRetries; i++ {
		err = RunMigrations(db)
		if err == nil {
			return nil
		}
		log.Printf("Migration %d failed: %v. Retrying...", i+1, err)
	}
	return fmt.Errorf("Failed after %d retries: %w", maxRetries, err)
}
