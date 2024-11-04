package database

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/nade-harlow/ecom-api/internal/config"
	"io/fs"
	"log"

	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/lib/pq"
)

//go:embed migrations/*
var migrationFiles embed.FS

func RunManualMigration() {
	dbConfig := config.AppConfig.Database
	dbURL := fmt.Sprintf("postgresql://%s:%s@%s:%v/%s?sslmode=disable",
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.DatabaseName)

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Use the embedded migration files
	migrations, err := fs.Sub(migrationFiles, "migrations")
	if err != nil {
		log.Fatal("Failed to create sub filesystem for migrations:", err)
	}

	d, err := iofs.New(migrations, ".")
	if err != nil {
		log.Fatal("Failed to create iofs driver:", err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal("Failed to create postgres driver:", err)
	}

	m, err := migrate.NewWithInstance("iofs", d, "postgres", driver)
	if err != nil {
		log.Fatal("Failed to create migrate instance:", err)
	}

	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatal("Failed to run migrations:", err)
	}

	log.Println("Completed DB migration")
}
