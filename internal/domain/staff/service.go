package staff

import (
	"errors"
	"net/http"
	"strconv"

	stafftype "github.com/mhvn092/movie-go/internal/domain/staff-type"
	"github.com/mhvn092/movie-go/internal/platform/web"
)

type StaffService struct {
	repo             *StaffRepository
	staffTypeService *stafftype.StaffTypeService
}

func NewStaffService(
	repo *StaffRepository,
	staffTypeService *stafftype.StaffTypeService,
) *StaffService {
	return &StaffService{repo: repo, staffTypeService: staffTypeService}
}

func (s *StaffService) GetAllPaginated(p web.PaginationParam) ([]StaffGetAllResponse, int, error) {
	return s.repo.getAllStaffPaginated(p)
}

func (s *StaffService) GetSearchResults(searchTerm string) ([]StaffGetAllResponse, error) {
	return s.repo.getSearchResults(searchTerm)
}

func (s *StaffService) CheckIfExists(id int) (bool, error) {
	return s.repo.checkIfExists(id)
}

func (s *StaffService) CheckCountOfExistingIds(ids []int) (bool, error) {
	return s.repo.checkCountOfExistingIds(ids)
}

func (s *StaffService) Insert(staff *Staff) (int, error) {
	exists, err := s.staffTypeService.CheckIfExists(staff.StaffTypeId)
	if err != nil {
		return 0, err
	}

	if !exists {
		return 0, errors.New(strconv.Itoa(http.StatusNotFound))
	}

	return s.repo.insert(staff)
}

func (s *StaffService) GetDetail(id int) (StaffGetDetailResponse, error) {
	return s.repo.getDetail(id)
}

func (s *StaffService) Edit(id int, staff *Staff) error {
	exists, err := s.staffTypeService.CheckIfExists(staff.StaffTypeId)
	if err != nil {
		return err
	}

	if !exists {
		return errors.New(strconv.Itoa(http.StatusNotFound))
	}

	return s.repo.edit(id, staff)
}

func (s *StaffService) Delete(id int) error {
	return s.repo.delete(id)
}
