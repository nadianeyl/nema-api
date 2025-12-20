package service

import (
	"github.com/nadianeyl/nema-api/internal/mailer"
	"github.com/nadianeyl/nema-api/internal/repository"
)

type Services struct {
	Users UserService
}

func NewServices(repositories repository.Repositories, m mailer.Mailer) Services {
	return Services{
		Users: NewUserService(repositories.Users, m),
	}
}
