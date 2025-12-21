package service

import (
	"time"

	"github.com/nadianeyl/nema-api/internal/jsonlog"
	"github.com/nadianeyl/nema-api/internal/repository"
)

type TokenService struct {
	UserRepo  repository.UserRepository
	TokenRepo repository.TokenRepository
	Logger    *jsonlog.Logger
}

func NewTokenService(userRepo repository.UserRepository, tokenRepo repository.TokenRepository, logger *jsonlog.Logger) TokenService {
	return TokenService{
		UserRepo:  userRepo,
		TokenRepo: tokenRepo,
		Logger:    logger,
	}
}

func (s *TokenService) CreateAuthToken(req *CreateAuthTokenRequest) (*TokenResponse, error) {
	user, err := s.UserRepo.GetByEmail(req.Email)
	if err != nil {
		return nil, err
	}

	match, err := user.Password.Matches(req.Password)
	if err != nil {
		return nil, err
	}

	if !match {
		return nil, repository.ErrInvalidAuthCredentials
	}

	token, err := s.TokenRepo.New(user.ID, 24*time.Hour, repository.ScopeAuthentication)
	if err != nil {
		return nil, err
	}

	res := &TokenResponse{
		Plaintext: token.Plaintext,
		Expiry:    token.Expiry,
	}

	return res, nil
}
