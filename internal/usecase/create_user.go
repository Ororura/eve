package usecase

import "eve/domain"

type CreateUserUseCase struct {
	repo UserRepository
}

func NewCreateUserUseCase(r UserRepository) *CreateUserUseCase {
	return &CreateUserUseCase{repo: r}
}

func (cu *CreateUserUseCase) Execute(user domain.User) error {
	return cu.repo.Save(user)
}
