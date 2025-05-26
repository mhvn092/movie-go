package repository

import "github.com/jackc/pgx/v5/pgxpool"

type BaseRepository struct {
	DB *pgxpool.Pool
}
