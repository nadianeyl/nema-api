package service

import (
	"context"
	"time"

	"github.com/nadianeyl/nema-api/internal/domain"
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

func (s *TokenService) CreateAuthToken(ctx context.Context, req *CreateAuthTokenRequest) (*TokenResponse, error) {
	user, err := s.UserRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	match, err := user.Password.Matches(req.Password)
	if err != nil {
		return nil, err
	}

	if !match {
		return nil, domain.ErrInvalidAuthCredentials
	}

	token, err := s.TokenRepo.New(ctx, user.ID, 24*time.Hour, domain.ScopeAuthentication)
	if err != nil {
		return nil, err
	}

	res := &TokenResponse{
		Plaintext: token.Plaintext,
		Expiry:    token.Expiry,
	}

	return res, nil
}
