package stafftype

import "github.com/mhvn092/movie-go/internal/platform/web"

type StaffTypeService struct {
	repo *StaffTypeRepository
}

func NewStaffTypeService(repo *StaffTypeRepository) *StaffTypeService {
	return &StaffTypeService{repo: repo}
}

func (s *StaffTypeService) GetAllPaginated(p web.PaginationParam) ([]StaffType, int, error) {
	return s.repo.getAllStaffTypesPaginated(p)
}

func (s *StaffTypeService) Insert(genre *StaffType) (int, error) {
	return s.repo.insert(genre)
}

func (s *StaffTypeService) Edit(id int, genre *StaffType) error {
	return s.repo.edit(id, genre)
}

func (s *StaffTypeService) Delete(id int) error {
	return s.repo.delete(id)
}

func (s *StaffTypeService) CheckIfExists(id int) (bool, error) {
	return s.repo.checkIfExistsById(id)
}
