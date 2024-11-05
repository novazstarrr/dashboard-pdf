// internal/config/config.go
package config

import (
    "os"
)

type Config struct {
    Port        string
    Environment string
    Database    DatabaseConfig
    JWT         JWTConfig
    File        FileConfig
}

type DatabaseConfig struct {
    DBPath string
}

type JWTConfig struct {
    Secret string
}

type FileConfig struct {
    UploadDir    string
    MaxSize      int64
    AllowedTypes []string
    BaseURL      string
}

func NewConfig() *Config {
    return &Config{
        Port:        getEnvOrDefault("PORT", "8080"),
        Environment: getEnvOrDefault("ENV", "development"),
        Database: DatabaseConfig{
            DBPath: getEnvOrDefault("DB_PATH", "./database.db"),
        },
        JWT: JWTConfig{
            Secret: getEnvOrDefault("JWT_SECRET", "your-default-secret-key"),
        },
        File: FileConfig{
            UploadDir:    getEnvOrDefault("UPLOAD_DIR", "./uploads"),
            MaxSize:      100 * 1024 * 1024, // 100MB default
            AllowedTypes: []string{
                "image/jpeg",
                "image/png",
                "image/gif",
                "application/pdf",
                "text/plain",
            },
            BaseURL: getEnvOrDefault("BACKEND_URL", "http://localhost:8080"),
        },
    }
}

func getEnvOrDefault(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}
