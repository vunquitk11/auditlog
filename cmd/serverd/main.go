package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/audit-log-service/cmd/common"
	"github.com/audit-log-service/cmd/serverd/router"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	// Init repository connection
	repoRegistry := common.InitRepository()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	corsOrigins := getCORSOrigins()
	r := router.New(ctx, corsOrigins, repoRegistry)

	appPort := getEnv("APP_PORT", "8080")
	srv := &http.Server{
		Addr:         ":" + appPort,
		Handler:      r.Handler(),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("Server starting on port %s", appPort)

	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("Server error: %v", err)
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getCORSOrigins() []string {
	origins := getEnv("CORS_ALLOWED_ORIGINS", "*")
	return strings.Split(origins, ",")
}
