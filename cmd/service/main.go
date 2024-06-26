package main

import (
	"github.com/jackc/pgx/v5/pgxpool"

	config "github.com/mhvn092/movie-go/internal"
	root "github.com/mhvn092/movie-go/internal/rest"
	"github.com/mhvn092/movie-go/internal/util"
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
	conn := util.InitDb()

	url, r := util.CreateServer()

	config.InitializeAppConfig(r, conn)

	root.InitializeRoutes()

	return conn, url, r
}
