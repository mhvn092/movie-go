package main

import (
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/mhvn092/movie-go/internal/platform/config"
	"github.com/mhvn092/movie-go/internal/platform/database"
	root "github.com/mhvn092/movie-go/internal/transport/http"
	"github.com/mhvn092/movie-go/pkg/exception"
	"github.com/mhvn092/movie-go/pkg/router"
)

func main() {
	conn, url, r := initialize()
	defer conn.Close()

	e := r.Serve(url)

	exception.ErrorExit(e, "server Creation Error")
}

func initialize() (*pgxpool.Pool, string, *router.Router) {
	conn := database.InitDb()

	url, r := root.CreateServer()

	config.InitializeAppConfig(r, conn)

	root.InitializeRoutes()

	return conn, url, r
}
