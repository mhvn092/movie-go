package main

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	config "github.com/mhvn092/movie-go/internal"
	root "github.com/mhvn092/movie-go/internal/rest"
	"github.com/mhvn092/movie-go/internal/util"
	"net/http"
)

func main() {
	conn, url, mux := initialize()
	defer conn.Close()

	e := http.ListenAndServe(url, mux)

	util.ErrorExit(e, "server Creation Error")
}

func initialize() (*pgxpool.Pool, string, *http.ServeMux) {
	conn := util.InitDb()

	url, mux := util.CreateServer()

	config.InitializeAppConfig(mux, conn)

	root.InitializeRoutes()

	return conn, url, mux
}
