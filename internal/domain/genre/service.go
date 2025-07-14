package genre

import "github.com/mhvn092/movie-go/internal/platform/web"

type GenreService struct {
	repo *GenreRepository
}

func NewGenreService(repo *GenreRepository) *GenreService {
	return &GenreService{repo: repo}
}

func (s *GenreService) GetAllPaginated(p web.PaginationParam) ([]Genre, int, error) {
	return s.repo.getAllGenresPaginated(p)
}

func (s *GenreService) Insert(genre *Genre) (int, error) {
	return s.repo.insert(genre)
}

func (s *GenreService) CheckIfExists(id int) (bool, error) {
	return s.repo.checkIfExistsById(id)
}

func (s *GenreService) Edit(id int, genre *Genre) error {
	return s.repo.edit(id, genre)
}

func (s *GenreService) Delete(id int) error {
	return s.repo.delete(id)
}
