package config

import (
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/mhvn092/movie-go/pkg/router"
)

type appConfigStruct struct {
	mux *router.Router
	db  *pgxpool.Pool
}

var appConfig *appConfigStruct

func GetDbPool() *pgxpool.Pool {
	return appConfig.db
}

func GetRouter() *router.Router {
	return appConfig.mux
}

func InitializeAppConfig(mux *router.Router, db *pgxpool.Pool) {
	appConfig = &appConfigStruct{mux, db}
}
