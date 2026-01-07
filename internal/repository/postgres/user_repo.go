package postgres

import (
	"eve/domain"

	"github.com/jmoiron/sqlx"
)

type UserRepo struct {
	db *sqlx.DB
}

func (u *UserRepo) Save(user domain.User) error {
	_, err := u.db.Exec("INSERT INTO users (email, password) VALUES ($1, $2)", user.Email, user.Password)
	return err
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
