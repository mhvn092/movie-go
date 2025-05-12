package genre

import (
	"github.com/jackc/pgx/v5/pgxpool"

	config "github.com/mhvn092/movie-go/internal"
	"github.com/mhvn092/movie-go/pkg/router"
)

var db *pgxpool.Pool

func Router() *router.Router {
	db = config.GetDbPool()
	r := router.NewRouter()
	r.GetWithPagination("/all", getAll)
	return r
}
