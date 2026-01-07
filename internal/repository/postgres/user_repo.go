package postgres

import (
	"eve/domain"

	"github.com/jmoiron/sqlx"
)

type UserRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) *UserRepo {
	return &UserRepo{db}
}

func (u *UserRepo) GetAll() ([]domain.User, error) {
	var users []domain.User
	err := u.db.Select(&users, `
		SELECT id,email,password,created_at FROM users
	`)

	return users, err
}
