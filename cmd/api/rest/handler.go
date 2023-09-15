package rest

import (
	"github.com/d-ashesss/mah-moneh/internal/accounts"
	"github.com/d-ashesss/mah-moneh/internal/auth"
	"github.com/d-ashesss/mah-moneh/internal/categories"
	"github.com/d-ashesss/mah-moneh/internal/spendings"
	"github.com/d-ashesss/mah-moneh/internal/transactions"
	"github.com/d-ashesss/mah-moneh/internal/users"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"log"
	"net/http"
	"regexp"
	"strings"
)

type handler struct {
	auth         *auth.Service
	accounts     *accounts.Service
	categories   *categories.Service
	transactions *transactions.Service
	spendings    *spendings.Service
}

func NewHandler(
	cfg *Config,
	auth *auth.Service,
	accounts *accounts.Service,
	categories *categories.Service,
	transactions *transactions.Service,
	spendings *spendings.Service,
) http.Handler {
	h := &handler{
		auth:         auth,
		accounts:     accounts,
		categories:   categories,
		transactions: transactions,
		spendings:    spendings,
	}

	r := gin.New()
	r.HandleMethodNotAllowed = true
	r.Use(gin.Logger(), gin.CustomRecoveryWithWriter(nil, h.handleRecovery))

	corsCfg := cors.DefaultConfig()
	if len(cfg.AllowedOrigins) > 0 {
		corsCfg.AllowOrigins = cfg.AllowedOrigins
		corsCfg.AllowCredentials = true
	} else {
		corsCfg.AllowAllOrigins = true
	}
	r.Use(cors.New(corsCfg))

	r.NoRoute(h.notFound)
	r.NoMethod(h.methodNotAllowed)

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("yearmonth", validateYearMonth)
	}

	r.GET("/", h.handleIndex)

	r.Use(h.authenticate)
	r.GET("/deep-vaults", h.handleIndex)

	r.POST("/accounts", h.handleAccountsCreate)
	r.GET("/accounts", h.handleAccountsList)
	r.PUT("/accounts/:uuid", h.handleAccountsUpdate)
	r.DELETE("/accounts/:uuid", h.handleAccountsDelete)

	r.PUT("/accounts/:uuid/amounts", h.handleAccountAmountSet)
	r.PUT("/accounts/:uuid/amounts/:month", h.handleAccountAmountSet)
	r.GET("/accounts/:uuid/amounts", h.handleAccountAmountGet)
	r.GET("/accounts/:uuid/amounts/:month", h.handleAccountAmountGet)

	r.POST("/categories", h.handleCategoriesCreate)
	r.GET("/categories", h.handleCategoriesList)
	r.DELETE("/categories/:uuid", h.handleCategoriesDelete)

	r.POST("/transactions", h.handleTransactionsCreate)
	r.GET("/transactions/:month", h.handleTransactionsList)
	r.DELETE("/transactions/:uuid", h.handleTransactionsDelete)

	r.GET("/spendings/:month", h.handleSpendingsGet)

	return r
}

func (h *handler) handleIndex(c *gin.Context) {
	c.String(http.StatusOK, "ok")
}

func (h *handler) authenticate(c *gin.Context) {
	token := c.GetHeader("Authorization")
	token, _ = strings.CutPrefix(token, "Bearer ")

	user, err := h.auth.AuthenticateUser(c, token)
	if err != nil {
		log.Printf("[HTTP] Unauthorized request: %s", err)
		c.AbortWithStatusJSON(http.StatusUnauthorized, NewErrorResponse(http.StatusText(http.StatusUnauthorized)))
		return
	}
	c.Set("user", user)
	c.Next()
}

func (h *handler) user(c *gin.Context) *users.User {
	user, ok := c.Get("user")
	if !ok {
		return nil
	}
	return user.(*users.User)
}

func validateYearMonth(fl validator.FieldLevel) bool {
	month, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}
	rx := regexp.MustCompile("^\\d{4}-\\d{2}$")
	return month == "" || rx.MatchString(month)
}
