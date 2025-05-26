package genre

import "github.com/mhvn092/movie-go/internal/platform/web"

type GenreService struct {
	repo *GenreRepository
}

func NewGenreService(repo *GenreRepository) *GenreService {
	return &GenreService{repo: repo}
}

func (s *GenreService) GetAllPaginated(p web.PaginationParam) ([]Genre, int, error) {
	return s.repo.GetAllGenresPaginated(p)
}
