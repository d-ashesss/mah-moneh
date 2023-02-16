package main

import (
	"fmt"
	"github.com/d-ashesss/mah-moneh/internal/accounts"
	"github.com/d-ashesss/mah-moneh/internal/api"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
)

var accountService *accounts.Service

func main() {
	log.SetFlags(0)

	dbHost := os.Getenv("POSTGRES_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}
	dbUser := os.Getenv("POSTGRES_USER")
	dbPwd := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")
	if dbName == "" {
		log.Fatalf("DB name not configured")
	}
	dsn := fmt.Sprintf("host=%s user=%s password=%s database=%s", dbHost, dbUser, dbPwd, dbName)
	dbPort := os.Getenv("POSTGRES_PORT")
	if dbPort != "" {
		dsn = fmt.Sprintf("%s port=%s", dsn, dbPort)
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: logger.Discard})
	if err != nil {
		log.Fatalf("Failed to connect to the DB: %s", err)
	}
	store := accounts.NewGormStore(db)
	accountService = accounts.NewService(store)
	if err := db.AutoMigrate(&accounts.Account{}); err != nil {
		log.Fatalf("Failed to run DB migration: %s", err)
	}

	apiService := api.NewService(accountService)

	cfg := LoadConfig()
	app := NewApp(cfg, apiService)
	app.Run()
}
