package service

import (
	"errors"
	"sync"
	"time"

	"github.com/nadianeyl/nema-api/internal/helper"
	"github.com/nadianeyl/nema-api/internal/jsonlog"
	"github.com/nadianeyl/nema-api/internal/mailer"
	"github.com/nadianeyl/nema-api/internal/repository"
)

type UserService struct {
	UserRepo  repository.UserRepository
	TokenRepo repository.TokenRepository
	Mailer    mailer.Mailer
	Logger    *jsonlog.Logger
}

func NewUserService(userRepo repository.UserRepository, tokenRepo repository.TokenRepository, m mailer.Mailer, logger *jsonlog.Logger) UserService {
	return UserService{
		UserRepo:  userRepo,
		TokenRepo: tokenRepo,
		Mailer:    m,
		Logger:    logger,
	}
}

func (s *UserService) Register(req *RegisterUserRequest, wg *sync.WaitGroup) (*UserResponse, error) {
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

	token, err := s.TokenRepo.New(user.ID, 3*24*time.Hour, repository.ScopeActivation)
	if err != nil {
		return nil, err
	}

	helper.Background(s.Logger, wg, func() {
		data := map[string]any{
			"activationToken": token.Plaintext,
		}

		err = s.Mailer.Send(user.Email, "user_welcome.tmpl", data)
		if err != nil {
			s.Logger.LogError(err, nil)
		}
	})

	res := &UserResponse{
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

func (s *UserService) Activate(req *ActivateUserRequest) (*UserResponse, error) {
	user, err := s.UserRepo.GetForToken(repository.ScopeActivation, req.TokenPlaintext)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrRecordNotFound):
			return nil, repository.ErrInvalidOrExpiredToken
		default:
			return nil, err
		}
	}

	user.Activated = true

	err = s.UserRepo.Update(user)
	if err != nil {
		return nil, err
	}

	err = s.TokenRepo.DeleteAllForUser(repository.ScopeActivation, user.ID)
	if err != nil {
		return nil, err
	}

	res := &UserResponse{
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
