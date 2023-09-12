package repo

import (
	"database/sql"

	"github.com/neglarken/clickhead/auth-ms/internal/model"
)

type SessionRepository struct {
	DB *sql.DB
}

func NewSessionRepository(db *sql.DB) *SessionRepository {
	return &SessionRepository{
		DB: db,
	}
}

func (r *SessionRepository) Upsert(user_id int, s *model.Session) error {
	return r.DB.QueryRow(
		"insert into session (user_id, refresh_token, exires_at) values ($1, $2, $3) on conflict (user_id) do update set (refresh_token, exires_at) = ($2, $3) where session.user_id = $1",
		user_id,
		s.RefreshToken,
		s.ExpiresAt).Err()
}
