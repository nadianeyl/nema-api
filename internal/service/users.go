package service

import (
	"github.com/nadianeyl/nema-api/internal/helper"
	"github.com/nadianeyl/nema-api/internal/jsonlog"
	"github.com/nadianeyl/nema-api/internal/mailer"
	"github.com/nadianeyl/nema-api/internal/repository"
)

type UserService struct {
	UserRepo repository.UserRepository
	Mailer   mailer.Mailer
	Logger   *jsonlog.Logger
}

func NewUserService(userRepo repository.UserRepository, m mailer.Mailer, logger *jsonlog.Logger) UserService {
	return UserService{
		UserRepo: userRepo,
		Mailer:   m,
		Logger:   logger,
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

	helper.Background(s.Logger, func() {
		err = s.Mailer.Send(user.Email, "user_welcome.tmpl", user)
		if err != nil {
			s.Logger.LogError(err, nil)
		}
	})

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
