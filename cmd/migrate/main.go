// Command migrate applies the embedded goose migrations to the database and
// exits. The bot also runs these on startup; this is for running them on their
// own (e.g. against a DB copy) without a bot token.
package main

import (
	"log/slog"
	"os"

	"weazyexe.dev/didntmaker/internal/database"
)

func main() {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "didntmaker.db"
	}

	db, err := database.Init(dbPath)
	if err != nil {
		slog.Error("migration failed", "error", err)
		os.Exit(1)
	}
	defer database.Close(db)

	slog.Info("migrations applied", "path", dbPath)
}
