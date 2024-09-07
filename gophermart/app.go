package gophermart

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"

	"github.com/novoseltcev/go-diploma-gofermart/gophermart/endpoints"
	"github.com/novoseltcev/go-diploma-gofermart/shared"

	"github.com/novoseltcev/go-diploma-gofermart/gophermart/auth"
)


type App struct {
	config Config
}


func NewApp(config Config) *App {
	return &App{
		config: config,
	}
}

func (app *App) Start() error {
	if app.config.DatabaseDsn == "" {
	    return errors.New("not specified DatabaseDsn from env or arg")
	}

	db, err := sqlx.Open("pgx", app.config.DatabaseDsn)
	if err != nil {
		return err
	}
	defer db.Close()

	log.Info("App is started")
	return http.ListenAndServe(app.config.Address, app.GetRouter(db))
}

func (app *App) GetRouter(db *sqlx.DB) http.Handler {
	r := gin.Default()
	r.Use(gzip.Gzip(gzip.DefaultCompression))

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, "pong")
	})

	jwtManager := auth.NewJWTManager(app.config.JwtSecret, time.Duration(app.config.JwtLifetime) * time.Hour * 24)
	uowPool := shared.NewUOWPool(db)

	r.POST("/api/user/login", endpoints.Login(uowPool, jwtManager))
	r.POST("/api/user/register", endpoints.Register(uowPool, jwtManager))


	userAPI := r.Group("/api/user")
	{
		userAPI.Use(auth.JWTMiddleware(jwtManager))

		userAPI.POST("/orders", endpoints.AddOrder(uowPool))
		userAPI.GET("/orders", endpoints.GetOrders(uowPool))

		userAPI.GET("/balance", endpoints.GetBalance(uowPool))
		userAPI.POST("/balance/withdraw", endpoints.Withdraw(uowPool))
		userAPI.GET("/withdrawals", endpoints.GetWithdrawals(uowPool))
	}

	return r
}
