package database

import (
	"log/slog"

	"weazyexe.dev/didntmaker/internal/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Init(dbPath string) (*gorm.DB, error) {
	slog.Info("opening database", "path", dbPath)

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		slog.Error("failed to open database", "error", err)
		return nil, err
	}

	slog.Info("running migrations")
	if err := db.AutoMigrate(&models.User{}); err != nil {
		slog.Error("failed to run migrations", "error", err)
		return nil, err
	}

	slog.Info("database ready")
	return db, nil
}
