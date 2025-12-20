package service

import "github.com/nadianeyl/nema-api/internal/repository"

type UserService struct {
	UserRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return UserService{
		UserRepo: userRepo,
	}
}

func (s *UserService) Register(req *RegisterUserRequest) (*RegisterUserResponse, error) {
	user := &repository.User{
		Name:                      req.Name,
		Email:                     req.Email,
		Activated:                 false,
		EmailNotificationsEnabled: true,
	}

	err := user.Password.Set(req.Password)
	if err != nil {
		return nil, err
	}

	err = s.UserRepo.Insert(user)
	if err != nil {
		return nil, err
	}

	res := &RegisterUserResponse{
		ID:                        user.ID,
		Name:                      user.Name,
		Email:                     user.Email,
		Activated:                 user.Activated,
		EmailNotificationsEnabled: user.EmailNotificationsEnabled,
		CreatedAt:                 user.CreatedAt,
		UpdatedAt:                 user.UpdatedAt,
	}

	return res, nil
}
