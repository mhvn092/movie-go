package config

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
)

type AppConfigStruct struct {
	Mux *http.ServeMux
	Db  *pgxpool.Pool
}

var AppConfig *AppConfigStruct

func InitializeAppConfig(mux *http.ServeMux, db *pgxpool.Pool) {
	AppConfig = &AppConfigStruct{mux, db}
}
