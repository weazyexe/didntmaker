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
	if err := db.AutoMigrate(&models.User{}, &models.Transaction{}); err != nil {
		slog.Error("failed to run migrations", "error", err)
		return nil, err
	}

	if err := migrateExistingBalances(db); err != nil {
		slog.Error("failed to migrate existing balances", "error", err)
	}

	return db, nil
}

func migrateExistingBalances(db *gorm.DB) error {
	var count int64
	db.Model(&models.Transaction{}).Where("type = ?", models.TransactionTypeMigration).Count(&count)
	if count > 0 {
		slog.Info("existing balances already migrated", "count", count)
		return nil
	}

	var users []models.User
	db.Where("balance != 0").Find(&users)

	if len(users) == 0 {
		slog.Info("no existing balances to migrate")
		return nil
	}

	// Wrap entire migration in a transaction for atomicity
	err := db.Transaction(func(dbTx *gorm.DB) error {
		for _, user := range users {
			tx := models.Transaction{
				ChatID:     user.ChatID,
				Type:       models.TransactionTypeMigration,
				SenderID:   0,
				ReceiverID: user.TelegramID,
				Amount:     user.Balance,
				CreatedAt:  user.CreatedAt,
			}
			if err := dbTx.Create(&tx).Error; err != nil {
				slog.Error("failed to create migration transaction",
					"chat_id", user.ChatID,
					"telegram_id", user.TelegramID,
					"balance", user.Balance,
					"error", err,
				)
				return err // rollback entire migration
			}
		}
		return nil
	})

	if err != nil {
		return err
	}

	slog.Info("migrated existing balances", "count", len(users))
	return nil
}

func Close(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	slog.Info("closing database connection")
	return sqlDB.Close()
}
