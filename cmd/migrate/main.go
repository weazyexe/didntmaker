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
