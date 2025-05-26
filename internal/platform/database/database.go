package database

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/mhvn092/movie-go/pkg/env"
	"github.com/mhvn092/movie-go/pkg/exception"
)

var pgOnce sync.Once

func createDb() *pgxpool.Pool {
	var conn *pgxpool.Pool
	var err error

	databaseUrl := env.GetEnv(env.DATABASE_URL)

	config, err := pgxpool.ParseConfig(databaseUrl)
	if err != nil {
		exception.ErrorExit(err, "Couldn't parse database URL")
	}
	config.MaxConns = 10
	config.MaxConnLifetime = time.Minute * 3
	config.MaxConnIdleTime = 10

	// Initialize connection with sync.Once to ensure it's done only once
	pgOnce.Do(func() {
		conn, err = pgxpool.NewWithConfig(context.Background(), config)
		if err != nil {
			exception.ErrorExit(err, "Couldn't connect to the database")
		}

		// Ping the database to ensure it's available
		err = conn.Ping(context.Background())
		if err != nil {
			exception.ErrorExit(err, "Database connection is unavailable")
		}
	})

	if conn == nil {
		exception.ErrorExit(fmt.Errorf("connection is nil"), "Database connection failed")
	}

	return conn
}

func InitDb() *pgxpool.Pool {
	env.ReadEnv()
	fmt.Println("Env File is Read")

	conn := createDb()

	fmt.Println("Db Connection Established")
	return conn
}
