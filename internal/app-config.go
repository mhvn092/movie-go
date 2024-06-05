package config

import (
	"github.com/jackc/pgx/v5"
	"net/http"
)

type AppConfigStruct struct {
	Mux *http.ServeMux
	Db  *pgx.Conn
}

var AppConfig *AppConfigStruct

func InitializeAppConfig(mux *http.ServeMux, db *pgx.Conn) {
	AppConfig = &AppConfigStruct{mux, db}
}
