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
	"github.com/d-ashesss/mah-moneh/internal/users"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/suite"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func NewAuthRequest(user *users.User, method, target string, body io.Reader) *http.Request {
	r := httptest.NewRequest(method, target, body)
	r.Header.Add("Authorization", "Bearer "+user.UUID.String())
	if body != nil {
		r.Header.Add("Content-Type", "application/json")
	}
	return r
}

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

	accountsService     *accounts.Service
	categoriesService   *categories.Service
	transactionsService *transactions.Service

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
	ts.accountsService = accounts.NewService(accountsStore)
	categoriesStore := categories.NewGormStore(db)
	ts.categoriesService = categories.NewService(categoriesStore)
	transactionsStore := transactions.NewGormStore(db)
	ts.transactionsService = transactions.NewService(transactionsStore)
	capitalService := capital.NewService(ts.accountsService)
	spendingsService := spendings.NewService(capitalService, ts.transactionsService, ts.categoriesService)

	if err := db.AutoMigrate(
		&accounts.Account{},
		&accounts.Amount{},
		&categories.Category{},
		&transactions.Transaction{},
	); err != nil {
		log.Fatalf("Failed to run DB migration: %s", err)
	}

	ts.handler = rest.NewHandler(ts.accountsService, ts.categoriesService, ts.transactionsService, spendingsService)
}

func (ts *RESTTestSuite) TestRest() {
	ts.Run("Index", ts.testIndex)
	ts.Run("Authorization", ts.testAuthorization)

	ts.Run("Errors", func() {
		ts.Run("Accounts", ts.testAccounts)
		ts.Run("Categories", ts.testCategories)
		ts.Run("Transactions", ts.testTransactions)
	})
}

func (ts *RESTTestSuite) testIndex() {
	request := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()
	ts.handler.ServeHTTP(rr, request)

	ts.Equal(200, rr.Code)
	ts.Equal("ok", rr.Body.String())
}

func (ts *RESTTestSuite) testAuthorization() {
	ts.Run("Unauthorized", func() {
		request := httptest.NewRequest("GET", "/deep-vaults", nil)
		rr := NewRecorder()
		ts.handler.ServeHTTP(rr, request)

		ts.Equal(http.StatusUnauthorized, rr.Code)
		r := new(ErrorResponse)
		rr.FromJSON(&r)
		ts.Equal(http.StatusText(http.StatusUnauthorized), r.Error)
	})

	ts.Run("Authorized", func() {
		auth := &users.User{UUID: uuid.Must(uuid.NewV4())}
		request := NewAuthRequest(auth, "GET", "/deep-vaults", nil)
		rr := NewRecorder()
		ts.handler.ServeHTTP(rr, request)

		ts.Equal(200, rr.Code)
		ts.Equal("ok", rr.Body.String())
	})
}

func TestREST(t *testing.T) {
	suite.Run(t, new(RESTTestSuite))
}
