package config

import (
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/mhvn092/movie-go/pkg/router"
)

type AppConfigStruct struct {
	Mux *router.Router
	Db  *pgxpool.Pool
}

var AppConfig *AppConfigStruct

func InitializeAppConfig(mux *router.Router, db *pgxpool.Pool) {
	AppConfig = &AppConfigStruct{mux, db}
}
