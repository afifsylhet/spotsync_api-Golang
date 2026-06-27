package auth

import (
	"errors"

	"github.com/afifsylhet/spotsync-api/internal/auth/dto"
	userdto "github.com/afifsylhet/spotsync-api/internal/user/dto"
	"github.com/afifsylhet/spotsync-api/internal/user"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService interface {
	Register(req dto.RegisterRequest) (*userdto.UserResponse, error)
	Login(req dto.LoginRequest) (*dto.LoginResponse, error)
}

type authService struct {
	userService user.UserService
}

func NewAuthService(userService user.UserService) AuthService {
	return &authService{
		userService: userService,
	}
}

func (s *authService) Register(req dto.RegisterRequest) (*userdto.UserResponse, error) {
	_, err := s.userService.FindByEmail(req.Email)
	if err == nil {
		return nil, errors.New("email already registered")
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		return nil, err
	}

	role := req.Role
	if role == "" {
		role = "driver"
	}

	u := &user.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     role,
	}

	if err := s.userService.CreateUser(u); err != nil {
		return nil, err
	}

	return &userdto.UserResponse{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Role:      u.Role,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}, nil
}

func (s *authService) Login(req dto.LoginRequest) (*dto.LoginResponse, error) {
	u, err := s.userService.FindByEmail(req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid email or password")
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(u.Password),
		[]byte(req.Password),
	); err != nil {
		return nil, errors.New("invalid email or password")
	}

	token, err := GenerateJWT(u.ID, u.Role)
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		Token: token,
		User: dto.LoginUserResponse{
			ID:    u.ID,
			Name:  u.Name,
			Email: u.Email,
			Role:  u.Role,
		},
	}, nil
}
