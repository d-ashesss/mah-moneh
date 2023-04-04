package main

import (
	"github.com/d-ashesss/mah-moneh/internal/accounts"
	"github.com/d-ashesss/mah-moneh/internal/api"
	"github.com/d-ashesss/mah-moneh/internal/capital"
	"github.com/d-ashesss/mah-moneh/internal/categories"
	"github.com/d-ashesss/mah-moneh/internal/datastore"
	"github.com/d-ashesss/mah-moneh/internal/spendings"
	"github.com/d-ashesss/mah-moneh/internal/transactions"
	"log"
)

func main() {
	log.SetFlags(0)

	dbCfg, err := datastore.NewConfig()
	if err != nil {
		log.Fatalf("Invalid database config: %s", err)
	}
	db, err := datastore.Open(dbCfg)
	if err != nil {
		log.Fatalf("Failed to connect to the DB: %s", err)
	}

	accountsStore := accounts.NewGormStore(db)
	accountsService := accounts.NewService(accountsStore)
	categoriesStore := categories.NewGormStore(db)
	categoriesService := categories.NewService(categoriesStore)
	transactionsStore := transactions.NewGormStore(db)
	transactionsService := transactions.NewService(transactionsStore)
	capitalService := capital.NewService(accountsService)
	spendingsService := spendings.NewService(capitalService, transactionsService, categoriesService)

	if err := db.AutoMigrate(
		&accounts.Account{},
		&accounts.Amount{},
		&categories.Category{},
		&transactions.Transaction{},
	); err != nil {
		log.Fatalf("Failed to run DB migration: %s", err)
	}

	apiService := api.NewService(accountsService, categoriesService, transactionsService, spendingsService)

	appCfg := NewConfig()
	app := NewApp(appCfg, apiService)
	app.Run()
}
