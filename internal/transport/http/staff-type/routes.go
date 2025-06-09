package stafftypehandler

import (
	stafftype "github.com/mhvn092/movie-go/internal/domain/staff-type"
	"github.com/mhvn092/movie-go/internal/platform/config"
	"github.com/mhvn092/movie-go/internal/platform/middleware"
	"github.com/mhvn092/movie-go/internal/platform/repository"
	"github.com/mhvn092/movie-go/pkg/router"
)

var service *stafftype.StaffTypeService

func initialize() {
	db := config.GetDbPool()
	staffTypeRepo := stafftype.NewStaffTypeRepository(&repository.BaseRepository{DB: db})
	service = stafftype.NewStaffTypeService(staffTypeRepo)
}

func Router() *router.Router {
	initialize()
	r := router.NewRouter()

	r.GetWithPagination("/all", getAll, middleware.AuthAdmin)
	r.Post("/create", insert, middleware.AuthAdmin)
	r.Put("/update/{id}", edit, middleware.AuthAdmin)
	r.Delete("/delete/{id}", delete, middleware.AuthAdmin)
	return r
}
