//go:build integration

package rest_test

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"github.com/d-ashesss/mah-moneh/cmd/api/rest"
	"github.com/d-ashesss/mah-moneh/internal/accounts"
	"github.com/d-ashesss/mah-moneh/internal/auth"
	"github.com/d-ashesss/mah-moneh/internal/capital"
	"github.com/d-ashesss/mah-moneh/internal/categories"
	"github.com/d-ashesss/mah-moneh/internal/datastore"
	"github.com/d-ashesss/mah-moneh/internal/spendings"
	"github.com/d-ashesss/mah-moneh/internal/transactions"
	"github.com/d-ashesss/mah-moneh/internal/users"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/stretchr/testify/suite"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

const TestKeyID = "test_key"

type RESTTestSuite struct {
	suite.Suite

	privKey jwk.Key

	accountsService     *accounts.Service
	categoriesService   *categories.Service
	transactionsService *transactions.Service

	handler http.Handler

	users struct {
		main    Auth
		control Auth
	}
	accounts struct {
		bank uuid.UUID
		cash uuid.UUID
		temp uuid.UUID
	}
	categories struct {
		income    uuid.UUID
		groceries uuid.UUID
		temp      uuid.UUID
	}
	transactions struct {
		temp uuid.UUID
	}
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

	rsaKey, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		log.Fatalf("Failed to generate private key: %s", err)
	}

	privKey, err := jwk.FromRaw(rsaKey)
	if err != nil {
		log.Fatalf("Failed to create public key: %s", err)
	}
	if err := privKey.Set(jwk.KeyIDKey, TestKeyID); err != nil {
		log.Fatalf("Failed to set key ID: %s", err)
	}
	if err := privKey.Set(jwk.AlgorithmKey, jwa.RS256); err != nil {
		log.Fatalf("Failed to set key ID: %s", err)
	}
	ts.privKey = privKey

	pubKey, err := jwk.FromRaw(rsaKey.Public())
	if err != nil {
		log.Fatalf("Failed to create public key: %s", err)
	}
	if err := pubKey.Set(jwk.KeyIDKey, TestKeyID); err != nil {
		log.Fatalf("Failed to set key ID: %s", err)
	}
	if err := pubKey.Set(jwk.AlgorithmKey, jwa.RS256); err != nil {
		log.Fatalf("Failed to set key ID: %s", err)
	}

	authCfg := auth.NewConfig()
	usersService := users.NewService()
	authService := auth.NewService(authCfg, usersService)
	if err := authService.AddKey(pubKey); err != nil {
		log.Fatalf("Failed to add test public key: %s", err)
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

	ts.handler = rest.NewHandler(
		authService,
		ts.accountsService,
		ts.categoriesService,
		ts.transactionsService,
		spendingsService,
	)

	ts.users.main = ts.NewAuth()
	ts.users.control = ts.NewAuth()
}

type Auth struct {
	UUID  uuid.UUID
	user  *users.User
	token string
}

func (ts *RESTTestSuite) NewAuth() Auth {
	UUID := uuid.Must(uuid.NewV4())
	t := jwt.New()
	if err := t.Set(jwt.SubjectKey, UUID.String()); err != nil {
		panic(err)
	}
	token, err := jwt.Sign(t, jwt.WithKey(ts.privKey.Algorithm(), ts.privKey))
	if err != nil {
		panic(err)
	}
	return Auth{
		UUID:  UUID,
		user:  &users.User{UUID: UUID},
		token: string(token),
	}
}

type Request struct {
	*http.Request
}

func NewRequest(method, target string, body io.Reader) *Request {
	r := &Request{
		Request: httptest.NewRequest(method, target, body),
	}
	if body != nil {
		r.Header.Add("Content-Type", "application/json")
	}
	return r
}

func (r *Request) WithAuth(a Auth) *Request {
	r.Header.Add("Authorization", "Bearer "+a.token)
	return r
}

func (ts *RESTTestSuite) Serve(request *Request) int {
	rr := httptest.NewRecorder()
	ts.handler.ServeHTTP(rr, request.Request)
	return rr.Code
}

func (ts *RESTTestSuite) ServeJSON(request *Request, response any) int {
	rr := httptest.NewRecorder()
	ts.handler.ServeHTTP(rr, request.Request)
	if err := json.Unmarshal(rr.Body.Bytes(), response); err != nil {
		panic(err)
	}
	return rr.Code
}

func (ts *RESTTestSuite) ServeString(request *Request) (int, string) {
	rr := httptest.NewRecorder()
	ts.handler.ServeHTTP(rr, request.Request)
	return rr.Code, rr.Body.String()
}

type RequestTest struct {
	Name   string
	Method string
	Target string
	Body   io.Reader
	Auth   Auth
	Code   int
}

type CreationTest struct {
	Name string
	Body io.Reader
	Ref  *uuid.UUID
}

type CreationTestResponse struct {
	UUID string `json:"uuid"`
}

type ErrorTest struct {
	Name   string
	Method string
	Target string
	Auth   Auth
	Body   io.Reader
	Code   int
	Error  string
}

type ErrorTestResponse struct {
	Error string `json:"error"`
}

type CountTest struct {
	Name   string
	Target string
	Auth   Auth
	Count  int
}

type JSONTest struct {
	Name     string
	Target   string
	Auth     Auth
	Expected string
}

func (ts *RESTTestSuite) testRequest(tt RequestTest) {
	ts.Run(tt.Name, func() {
		request := NewRequest(tt.Method, tt.Target, tt.Body).WithAuth(tt.Auth)
		code := ts.Serve(request)

		ts.Equal(tt.Code, code)
	})
}

func (ts *RESTTestSuite) testCreate(tt CreationTest, target string) {
	ts.Run(tt.Name, func() {
		request := NewRequest("POST", target, tt.Body).WithAuth(ts.users.main)
		response := new(CreationTestResponse)
		code := ts.ServeJSON(request, response)

		ts.Equal(http.StatusCreated, code)
		ts.Require().NotEmptyf(response.UUID, "Received invalid UUID value in response")
		if tt.Ref != nil {
			*tt.Ref = uuid.Must(uuid.FromString(response.UUID))
		}
	})
}

func (ts *RESTTestSuite) testError(tt ErrorTest) {
	ts.Run(tt.Name, func() {
		request := NewRequest(tt.Method, tt.Target, tt.Body).WithAuth(tt.Auth)
		response := new(ErrorTestResponse)
		code := ts.ServeJSON(request, response)

		ts.Equal(tt.Code, code)
		ts.Equal(tt.Error, response.Error)
	})
}

func (ts *RESTTestSuite) testCount(tt CountTest) {
	ts.Run(tt.Name, func() {
		request := NewRequest("GET", tt.Target, nil).WithAuth(tt.Auth)
		response := make([]map[string]any, 0)
		code := ts.ServeJSON(request, &response)

		ts.Equal(http.StatusOK, code)
		ts.Len(response, tt.Count)
	})
}

func (ts *RESTTestSuite) testJSON(tt JSONTest) {
	ts.Run(tt.Name, func() {
		request := NewRequest("GET", tt.Target, nil).WithAuth(tt.Auth)
		code, response := ts.ServeString(request)

		ts.Equal(http.StatusOK, code)
		ts.JSONEq(tt.Expected, response)
	})
}

func (ts *RESTTestSuite) TestREST() {
	ts.Run("Index", ts.testIndex)
	ts.Run("Authorization", ts.testAuthorization)

	ts.Run("Errors", func() {
		ts.Run("Accounts", ts.testAccountsErrors)
		ts.Run("Categories", ts.testCategoriesErrors)
		ts.Run("Transactions", ts.testTransactions)
	})

	ts.Run("Create", func() {
		ts.Run("Accounts", ts.testCreateAccounts)
		ts.Run("Categories", ts.testCreateCategories)
		ts.Run("Transactions", ts.testCreateTransactions)
	})

	ts.Run("Delete", func() {
		ts.Run("Accounts", ts.testDeleteAccounts)
		ts.Run("Categories", ts.testDeleteCategories)
		ts.Run("Transactions", ts.testDeleteTransactions)
	})

	ts.Run("Get", func() {
		ts.Run("Accounts", ts.testGetAccounts)
		ts.Run("Categories", ts.testGetCategories)
		ts.Run("Transactions", ts.testGetTransactions)
		ts.Run("Spendings", ts.testGetSpendings)
	})
}

func (ts *RESTTestSuite) testIndex() {
	request := NewRequest("GET", "/", nil)
	code, response := ts.ServeString(request)

	ts.Equal(http.StatusOK, code)
	ts.Equal("ok", response)
}

func (ts *RESTTestSuite) testAuthorization() {
	ts.Run("Unauthorized", func() {
		request := NewRequest("GET", "/deep-vaults", nil)
		response := new(ErrorTestResponse)
		code := ts.ServeJSON(request, response)

		ts.Equal(http.StatusUnauthorized, code)
		ts.Equal(http.StatusText(http.StatusUnauthorized), response.Error)
	})

	ts.Run("Authorized", func() {
		request := NewRequest("GET", "/deep-vaults", nil).WithAuth(ts.NewAuth())
		code, response := ts.ServeString(request)

		ts.Equal(http.StatusOK, code)
		ts.Equal("ok", response)
	})
}

func TestREST(t *testing.T) {
	suite.Run(t, new(RESTTestSuite))
}
