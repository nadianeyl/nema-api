package service

import (
	"github.com/nadianeyl/nema-api/internal/jsonlog"
	"github.com/nadianeyl/nema-api/internal/mailer"
	"github.com/nadianeyl/nema-api/internal/repository"
)

type Services struct {
	Users      UserService
	Tokens     TokenService
	Accounts   AccountService
	Categories CategoryService
}

func NewServices(repositories repository.Repositories, m mailer.Mailer, logger *jsonlog.Logger) Services {
	return Services{
		Users:      NewUserService(repositories.Users, repositories.Tokens, m, logger),
		Tokens:     NewTokenService(repositories.Users, repositories.Tokens, logger),
		Accounts:   NewAccountService(repositories.Accounts),
		Categories: NewCategoryService(repositories.Categories),
	}
}
