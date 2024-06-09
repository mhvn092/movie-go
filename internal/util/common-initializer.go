package util

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
	"time"
)

func InitDb() *pgxpool.Pool {
	ReadEnv()
	fmt.Printf("Env File is Read\n")

	databaseUrl := GetEnv(EnvKeys.DATABASE_URL)

	config, err := pgxpool.ParseConfig(databaseUrl)
	if err != nil {
		ErrorExit(err, "Couldn't parse database url")
	}
	config.MaxConns = 10
	config.MaxConnLifetime = time.Minute * 3
	config.MaxConnIdleTime = 10
	conn, e := pgxpool.NewWithConfig(context.Background(), config)

	if e != nil {
		ErrorExit(e, "Unable to connect to database")
	}
	fmt.Printf("Db Connection Established\n")
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
