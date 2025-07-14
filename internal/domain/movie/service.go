package movie

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/mhvn092/movie-go/internal/domain/genre"
	"github.com/mhvn092/movie-go/internal/domain/staff"
	stafftype "github.com/mhvn092/movie-go/internal/domain/staff-type"
	"github.com/mhvn092/movie-go/internal/platform/web"
)

type MovieService struct {
	repo             *MovieRepository
	staffTypeService *stafftype.StaffTypeService
	staffService     *staff.StaffService
	genreService     *genre.GenreService
}

func NewMovieService(
	repo *MovieRepository,
	staffTypeService *stafftype.StaffTypeService,
	staffService *staff.StaffService,
	genreService *genre.GenreService,
) *MovieService {
	return &MovieService{
		repo:             repo,
		staffTypeService: staffTypeService,
		staffService:     staffService,
		genreService:     genreService,
	}
}

func (s *MovieService) GetAllPaginated(p web.PaginationParam) ([]MovieGetAllResponse, int, error) {
	return s.repo.getAllMoviePaginated(p)
}

func (s *MovieService) GetSearchResults(searchTerm string) ([]MovieGetAllResponse, error) {
	return s.repo.getSearchResults(searchTerm)
}

func (s *MovieService) Insert(payload *MovieUpsertPayload) (int, error) {
	if err := s.validateUpsertPayload(payload); err != nil {
		return 0, err
	}

	return s.repo.insert(payload)
}

func (s *MovieService) GetDetail(id int) (MovieGetDetailResponse, error) {
	return s.repo.getDetail(id)
}

func (s *MovieService) Edit(id int, payload *MovieUpsertPayload) error {
	exists, err := s.repo.checkIfExists(id)
	if err != nil {
		return err
	}

	if !exists {
		return errors.New(strconv.Itoa(http.StatusNotFound))
	}

	if err = s.validateUpsertPayload(payload); err != nil {
		return err
	}

	return s.repo.edit(id, payload)
}

func (s *MovieService) Delete(id int) error {
	exists, err := s.repo.checkIfExists(id)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New(strconv.Itoa(http.StatusNotFound))
	}
	return s.repo.delete(id)
}

func (s *MovieService) validateUpsertPayload(payload *MovieUpsertPayload) error {
	exists, err := s.genreService.CheckIfExists(payload.GenreId)
	if err != nil {
		return err
	}

	if !exists {
		return errors.New(strconv.Itoa(http.StatusNotFound))
	}

	staffIds, staffTypeIds := collectUpsertUniqueIds(payload)

	if err := s.validateUpsertPayloadIds(staffIds, s.staffService.CheckCountOfExistingIds, "staff"); err != nil {
		return err
	}

	if err := s.validateUpsertPayloadIds(staffTypeIds, s.staffTypeService.CheckCountOfExistingIds, "staff type"); err != nil {
		return err
	}

	return nil
}

func (s *MovieService) validateUpsertPayloadIds(
	ids []int,
	checkFunc func([]int) (bool, error),
	resource string,
) error {
	if len(ids) == 0 {
		return nil
	}
	exists, err := checkFunc(ids)
	if err != nil {
		return fmt.Errorf("failed to check %s IDs: %w", resource, err)
	}
	if !exists {
		return errors.New(strconv.Itoa(http.StatusNotFound))
	}
	return nil
}

func collectUpsertUniqueIds(payload *MovieUpsertPayload) ([]int, []int) {
	staffIds := make(map[int]bool)
	staffTypeIds := make(map[int]bool)

	staffIds[payload.DirectorId] = true

	for _, item := range payload.Staffs {
		staffIds[item.StaffId] = true
		staffTypeIds[item.StaffTypeId] = true
	}

	return toSlice(staffIds), toSlice(staffTypeIds)
}

func toSlice(set map[int]bool) []int {
	result := make([]int, 0, len(set))
	for id := range set {
		result = append(result, id)
	}
	return result
}
