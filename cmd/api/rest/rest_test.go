//go:build integration

package rest_test

import (
	"encoding/json"
	"github.com/d-ashesss/mah-moneh/cmd/api/rest"
	"github.com/d-ashesss/mah-moneh/internal/accounts"
	"github.com/d-ashesss/mah-moneh/internal/capital"
	"github.com/d-ashesss/mah-moneh/internal/categories"
	"github.com/d-ashesss/mah-moneh/internal/datastore"
	"github.com/d-ashesss/mah-moneh/internal/spendings"
	"github.com/d-ashesss/mah-moneh/internal/transactions"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

type ResponseRecorder struct {
	*httptest.ResponseRecorder
}

func NewRecorder() *ResponseRecorder {
	return &ResponseRecorder{
		ResponseRecorder: httptest.NewRecorder(),
	}
}

func (rr *ResponseRecorder) FromJSON(v any) {
	if err := json.Unmarshal(rr.Body.Bytes(), v); err != nil {
		panic(err)
	}
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type RESTTestSuite struct {
	suite.Suite

	handler http.Handler
}

func (ts *RESTTestSuite) SetupSuite() {
	gin.SetMode(gin.ReleaseMode)

	dbCfg, err := datastore.NewConfig()
	if err != nil {
		log.Fatalf("Invalid database config: %s", err)
	}
	dbCfg.TablePrefix = "rest_test_"
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

	ts.handler = rest.NewHandler(accountsService, categoriesService, transactionsService, spendingsService)
}

func (ts *RESTTestSuite) TestRest() {
	ts.Run("index", ts.testIndex)
}

func (ts *RESTTestSuite) testIndex() {
	request := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()
	ts.handler.ServeHTTP(rr, request)

	ts.Equal(200, rr.Code)
	ts.Equal("ok", rr.Body.String())
}

func TestREST(t *testing.T) {
	suite.Run(t, new(RESTTestSuite))
}
