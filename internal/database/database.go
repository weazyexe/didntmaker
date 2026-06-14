package database

import (
	"database/sql"
	"log/slog"

	"github.com/pressly/goose/v3"
	_ "modernc.org/sqlite"

	"weazyexe.dev/didntmaker/db/migrations"
)

func Init(dbPath string) (*sql.DB, error) {
	slog.Info("opening database", "path", dbPath)

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		slog.Error("failed to open database", "error", err)
		return nil, err
	}

	if err := db.Ping(); err != nil {
		slog.Error("failed to connect to database", "error", err)
		return nil, err
	}

	// SQLite serializes writes; a single connection avoids "database is locked".
	db.SetMaxOpenConns(1)

	if err := runMigrations(db); err != nil {
		slog.Error("failed to run migrations", "error", err)
		return nil, err
	}

	return db, nil
}

func runMigrations(db *sql.DB) error {
	slog.Info("running migrations")

	goose.SetBaseFS(migrations.FS)
	if err := goose.SetDialect("sqlite3"); err != nil {
		return err
	}
	return goose.Up(db, ".")
}

func Close(db *sql.DB) error {
	slog.Info("closing database connection")
	return db.Close()
}
