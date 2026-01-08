package usecase

import "eve/domain"

type UserRepository interface {
	Save(domain.User) error
	GetAll() ([]domain.User, error)
}

type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(password, hash string) (bool, error)
}
