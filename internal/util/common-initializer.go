package util

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mhvn092/movie-go/pkg/env"
	"github.com/mhvn092/movie-go/pkg/exception"
	"github.com/mhvn092/movie-go/pkg/router"
	"sync"
	"time"
)

var (
	pgOnce sync.Once
)

func createDb() *pgxpool.Pool {
	var conn *pgxpool.Pool
	var err error

	databaseUrl := env.GetEnv(env.DATABASE_URL)

	config, err := pgxpool.ParseConfig(databaseUrl)
	if err != nil {
		exception.ErrorExit(err, "Couldn't parse database url")
	}
	config.MaxConns = 10
	config.MaxConnLifetime = time.Minute * 3
	config.MaxConnIdleTime = 10

	pgOnce.Do(func() {
		conn, err = pgxpool.NewWithConfig(context.Background(), config)
		if err != nil {
			exception.ErrorExit(err, "Couldn't parse database url")
		}
	})

	return conn
}

func InitDb() *pgxpool.Pool {
	env.ReadEnv()
	fmt.Println("Env File is Read")

	conn := createDb()

	fmt.Println("Db Connection Established")
	return conn
}

func CreateServer() (string, *router.Router) {
	r := router.NewRouter()

	host := env.GetEnv(env.HOST)
	port := env.GetEnv(env.PORT)

	url := host + ":" + port
	fmt.Println("Listening on " + url)

	return url, r
}
