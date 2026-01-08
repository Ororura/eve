package usecase

import "eve/domain"

type CreateUserUseCase struct {
	repo           UserRepository
	passwordHasher PasswordHasher
}

func NewCreateUserUseCase(r UserRepository, h PasswordHasher) *CreateUserUseCase {
	return &CreateUserUseCase{repo: r, passwordHasher: h}
}

func (cu *CreateUserUseCase) Execute(user domain.User) error {
	hashedPassword, err := cu.passwordHasher.Hash(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword
	return cu.repo.Save(user)
}
