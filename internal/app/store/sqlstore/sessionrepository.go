package sqlstore

import (
	"database/sql"
	"redditclone/internal/app/model"
	"redditclone/internal/app/store"
)

type SessionRepository struct {
	store *Store
}

func (r *SessionRepository) Create(s *model.Session) error {
	return r.store.db.QueryRow(
		"INSERT INTO sessions (user_id,request_id,token) VALUES ((SELECT id from users WHERE id=$1),$2,$3)",
		s.UserID,
		s.RequestID,
		s.Token,
	).Err()
}

func (r *SessionRepository) FindByToken(token string) (*model.Session, error) {
	s := &model.Session{}
	if err := r.store.db.QueryRow(
		"SELECT user_id,request_id,token FROM sessions WHERE token = $1",
		token,
	).Scan(
		&s.UserID,
		&s.RequestID,
		&s.Token,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}
	return s, nil
}

func (r *SessionRepository) FindByUserID(userID string) (*model.Session, error) {
	s := &model.Session{}
	if err := r.store.db.QueryRow(
		"SELECT user_id,request_id,token FROM sessions WHERE user_id = $1",
		userID,
	).Scan(
		&s.UserID,
		&s.RequestID,
		&s.Token,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}
	return s, nil
}

func (r *SessionRepository) Delete(userID string) error {
	if err := r.store.db.QueryRow(
		"DELETE FROM sessions WHERE user_id = $1",
		userID,
	).Err(); err != nil {
		return err
	}
	return nil
}
