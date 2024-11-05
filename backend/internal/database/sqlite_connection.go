// internal/database/sqlite_connection.go
package database

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"tech-test/backend/internal/config"
	"tech-test/backend/internal/domain"
)

func SetupDB(cfg config.DatabaseConfig) (*gorm.DB, error) {
	gormLogger := logger.New(
		log.Default(),
		logger.Config{
			SlowThreshold:             time.Second,     
			LogLevel:                  logger.Info,     
			IgnoreRecordNotFoundError: true,           
			Colorful:                  true,            
		},
	)

	gormConfig := &gorm.Config{
		Logger: gormLogger,
		NowFunc: func() time.Time {
			return time.Now().UTC() 
		},
	}

	db, err := gorm.Open(sqlite.Open(cfg.DBPath), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	err = db.Transaction(func(tx *gorm.DB) error {
		if err := tx.AutoMigrate(&domain.File{}); err != nil {
			return fmt.Errorf("failed to migrate schema: %w", err)
		}

		if err := tx.Exec("UPDATE files SET content_type = mime_type WHERE content_type IS NULL").Error; err != nil {
			return fmt.Errorf("failed to populate content_type: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("migration failed: %w", err)
	}

	return db, nil
}


func CloseDB(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("failed to close database connection: %w", err)
	}

	return nil
}
