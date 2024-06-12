package util

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
	"sync"
	"time"
)

var (
	pgOnce sync.Once
)

func createDb() *pgxpool.Pool {
	var conn *pgxpool.Pool
	var err error

	databaseUrl := GetEnv(EnvKeys.DATABASE_URL)

	config, err := pgxpool.ParseConfig(databaseUrl)
	if err != nil {
		ErrorExit(err, "Couldn't parse database url")
	}
	config.MaxConns = 10
	config.MaxConnLifetime = time.Minute * 3
	config.MaxConnIdleTime = 10

	pgOnce.Do(func() {
		conn, err = pgxpool.NewWithConfig(context.Background(), config)
		if err != nil {
			ErrorExit(err, "Couldn't parse database url")
		}
	})

	return conn
}

func InitDb() *pgxpool.Pool {
	ReadEnv()
	fmt.Println("Env File is Read")

	conn := createDb()

	fmt.Println("Db Connection Established")
	return conn
}

func CreateServer() (string, *http.ServeMux) {
	mux := http.NewServeMux()

	host := GetEnv(EnvKeys.HOST)
	port := GetEnv(EnvKeys.PORT)

	url := host + ":" + port
	fmt.Println("Listening on " + url)

	return url, mux
}
