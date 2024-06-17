package gophermart

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	chi_middleware "github.com/go-chi/chi/v5/middleware"
	log "github.com/sirupsen/logrus"
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
	    return errors.New("Not specified DatabaseDsn from env or arg")
	}
	// storage, err := pg.New(app.config.DatabaseDsn)
	// if err != nil {
	// 	return nil
	// }
	// defer storage.Close()

	log.Info("App is started")
	return http.ListenAndServe(app.config.Address, app.GetRouter())
}

func (app *App) GetRouter() http.Handler {
	r := chi.NewRouter()
	r.Use(chi_middleware.Logger)
	r.Use(chi_middleware.Compress(4))
	r.Use(chi_middleware.Heartbeat("/ping"))

	return r
}
