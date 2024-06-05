package util

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"net/http"
)

func InitDb() *pgx.Conn {
	ReadEnv()
	fmt.Printf("Env File is Read\n")

	databaseUrl := GetEnv(EnvKeys.DATABASE_URL)
	conn, e := pgx.Connect(context.Background(), databaseUrl)
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
