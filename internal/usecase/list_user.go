package usecase

import "eve/domain"

type GetUserUseCase struct {
	repo UserRepository
}

func NewGetUserUseCase(r UserRepository) *GetUserUseCase {
	return &GetUserUseCase{repo: r}
}

func (cu *GetUserUseCase) Execute() ([]domain.User, error) {
	return cu.repo.GetAll()
}
