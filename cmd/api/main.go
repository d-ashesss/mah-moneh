package main

import (
	"github.com/d-ashesss/mah-moneh/cmd/api/rest"
	"github.com/d-ashesss/mah-moneh/internal/accounts"
	"github.com/d-ashesss/mah-moneh/internal/auth"
	"github.com/d-ashesss/mah-moneh/internal/capital"
	"github.com/d-ashesss/mah-moneh/internal/categories"
	"github.com/d-ashesss/mah-moneh/internal/datastore"
	"github.com/d-ashesss/mah-moneh/internal/spendings"
	"github.com/d-ashesss/mah-moneh/internal/transactions"
	"github.com/d-ashesss/mah-moneh/internal/users"
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

	authCfg := auth.NewConfig()
	usersService := users.NewService()
	authService := auth.NewService(authCfg, usersService)
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

	handlerCfg := rest.NewConfig()
	handler := rest.NewHandler(
		handlerCfg,
		authService,
		accountsService,
		categoriesService,
		transactionsService,
		spendingsService,
	)

	appCfg := NewConfig()
	app := NewApp(appCfg, handler)
	app.Run()
}
