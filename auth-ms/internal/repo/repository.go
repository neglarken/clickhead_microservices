package repo

import (
	"database/sql"

	"github.com/neglarken/clickhead/auth-ms/internal/model"
)

type User interface {
	Create(u *model.User) error
	FindByLogin(login string) (*model.User, error)
}

type Session interface {
	Create(user_id int, s *model.Session)
}

type Repository struct {
	User    *UserRepository
	Session *SessionRepository
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		User:    NewUserRepository(db),
		Session: NewSessionRepository(db),
	}
}
