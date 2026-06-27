package user

type UserService interface {
	CreateUser(user *User) error
	FindByEmail(email string) (*User, error)
	FindByID(id uint) (*User, error)
}

type userService struct {
	userRepo UserRepository
}

func NewUserService(userRepo UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) CreateUser(user *User) error {
	return s.userRepo.CreateUser(user)
}

func (s *userService) FindByEmail(email string) (*User, error) {
	return s.userRepo.FindByEmail(email)
}

func (s *userService) FindByID(id uint) (*User, error) {
	return s.userRepo.FindByID(id)
}
