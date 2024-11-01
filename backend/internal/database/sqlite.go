// internal/database/sqlite.go
package database

import (
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
    "tech-test/backend/internal/domain"
)

func SetupDatabase() (*gorm.DB, error) {
    db, err := gorm.Open(sqlite.Open("./data/files.db"), &gorm.Config{})
    if err != nil {
        return nil, err
    }

    err = db.AutoMigrate(&domain.File{})
    if err != nil {
        return nil, err
    }

    return db, nil
}
