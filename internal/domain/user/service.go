package user

import (
	"errors"

	"github.com/mhvn092/movie-go/internal/platform/security"
)

type UserService struct {
	repo *UserRepository
}

func NewUserService(repo *UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Register(u *User, isAdmin bool) error {
	if isAdmin {
		u.Role = UserRole.ADMIN
	}
	return s.repo.RegisterUser(u)
}

func (s *UserService) Login(loginDto *LoginDto) (*User, error) {
	user, err := s.repo.CheckUser(loginDto)
	if err != nil {
		return nil, err
	}

	if err := security.ComparePasswords(user.Password, loginDto.Password); err != nil {
		return nil, errors.New("email or password is incorrect")
	}

	return user, nil
}

func (s *UserService) GenerateToken(u *User) (string, error) {
	return security.CreateToken(security.UserTokenData{
		ID:    u.Id,
		Email: u.Email,
		Role:  string(u.Role),
	})
}
