package common

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/audit-log-service/internal/repository"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// InitRepository reads DB config from environment variables, initializes the DB connection, and returns the repository registry.
func InitRepository() repository.Registry {
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, user, pass, dbname, port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to DB: %v", err))
	}

	// Set connection pool parameters from env
	sqlDB, err := db.DB()
	if err != nil {
		panic(fmt.Sprintf("Failed to get sql.DB from gorm: %v", err))
	}

	maxOpenConns := getEnvInt("DB_MAX_OPEN_CONNS", 10)
	maxIdleConns := getEnvInt("DB_MAX_IDLE_CONNS", 5)
	connMaxLifetime := getEnvDuration("DB_CONN_MAX_LIFETIME", 30*time.Minute)

	sqlDB.SetMaxOpenConns(maxOpenConns)
	sqlDB.SetMaxIdleConns(maxIdleConns)
	sqlDB.SetConnMaxLifetime(connMaxLifetime)

	return repository.New(db)
}

// getEnvInt reads an int from env or returns default
func getEnvInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return def
}

// getEnvDuration reads a duration from env or returns default
func getEnvDuration(key string, def time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		d, err := time.ParseDuration(v)
		if err == nil {
			return d
		}
	}
	return def
}
