package service

import (
"errors"
"os"
"time"
"github.com/golang-jwt/jwt/v5"
"github.com/afifsylhet/spotsync-api/dto"
"github.com/afifsylhet/spotsync-api/models"
"github.com/afifsylhet/spotsync-api/repository"
"golang.org/x/crypto/bcrypt"
"gorm.io/gorm"
)


type AuthService interface {
Register(req dto.RegisterRequest) (*dto.UserResponse, error)
Login(req dto.LoginRequest) (*dto.LoginResponse, error)
}
type authService struct {
userRepo repository.UserRepository
}
func NewAuthService(userRepo repository.UserRepository) AuthService {
return &authService{userRepo: userRepo}
}
func (s *authService) Register(req dto.RegisterRequest) (*dto.UserResponse, err
// Check if email already exists
_, err := s.userRepo.FindByEmail(req.Email)
if err == nil {
return nil, errors.New("email already registered")
}
// Hash the password (cost 12)
hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12
if err != nil {
return nil, err
}
role := req.Role
if role == "" {
        role = "driver" // default role
}
    user := &models.User{
        Name:     req.Name,
        Email:    req.Email,
        Password: string(hashedPassword),
        Role:     role,
}
if err := s.userRepo.CreateUser(user); err != nil {
return nil, err
}
return &dto.UserResponse{
        ID:        user.ID,
        Name:      user.Name,
        Email:     user.Email,
        Role:      user.Role,
        CreatedAt: user.CreatedAt,
        UpdatedAt: user.UpdatedAt,
}, nil
}
func (s *authService) Login(req dto.LoginRequest) (*dto.LoginResponse, error) {
// Find user by email
    user, err := s.userRepo.FindByEmail(req.Email)
if err != nil {
if errors.Is(err, gorm.ErrRecordNotFound) {
return nil, errors.New("invalid email or password")
}
return nil, err
}
// Compare bcrypt hash
if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.P
return nil, errors.New("invalid email or password")
}
// Generate JWT
    token, err := generateJWT(user.ID, user.Role)
if err != nil {
return nil, err
}
service/zone_service.go
return &dto.LoginResponse{
        Token: token,
        User: dto.UserResponse{
            ID:    user.ID,
            Name:  user.Name,
            Email: user.Email,
            Role:  user.Role,
},
}, nil
}
// generateJWT creates a signed JWT token with user id and role
func generateJWT(userID uint, role string) (string, error) {
    secret := os.Getenv("JWT_SECRET")
    claims := jwt.MapClaims{
"user_id": userID,
"role":    role,
"exp":     time.Now().Add(24 * time.Hour).Unix(), // expires in 24 hour
}
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
return token.SignedString([]byte(secret))
}