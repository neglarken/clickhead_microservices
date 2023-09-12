package repo

import (
	"database/sql"

	"github.com/neglarken/clickhead/auth-ms/internal/model"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		DB: db,
	}
}

func (r *UserRepository) Create(u *model.User) error {
	return r.DB.QueryRow(
		"insert into users (login, hashed_password) VALUES ($1, $2) RETURNING id",
		u.Login,
		u.Password,
	).Scan(&u.Id)
}

func (r *UserRepository) FindByLogin(login string) (*model.User, error) {
	u := &model.User{}
	if err := r.DB.QueryRow("SELECT * FROM users WHERE login = $1", login).Scan(
		&u.Id,
		&u.Login,
		&u.Password,
	); err != nil {
		return nil, err
	}
	return u, nil
}

func (r *UserRepository) FindById(id int) (*model.User, error) {
	u := &model.User{}
	if err := r.DB.QueryRow(
		"SELECT id, login, hashed_password FROM users WHERE id = $1",
		id,
	).Scan(
		&u.Id,
		&u.Login,
		&u.Password,
	); err != nil {
		return nil, err
	}

	return u, nil
}
