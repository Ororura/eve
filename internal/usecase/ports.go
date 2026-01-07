package usecase

import "eve/domain"

type UserRepository interface {
	Save(domain.User) error
	GetAll() ([]domain.User, error)
}
