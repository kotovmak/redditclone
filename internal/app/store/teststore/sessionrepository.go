package teststore

import (
	"redditclone/internal/app/model"
	"redditclone/internal/app/store"
)

type SessionRepository struct {
	store    *Store
	sessions map[string]*model.Session
}

func (r *SessionRepository) Create(s *model.Session) error {
	r.sessions[s.Token] = s
	return nil
}

func (r *SessionRepository) FindByToken(token string) (*model.Session, error) {
	u, ok := r.sessions[token]
	if !ok {
		return nil, store.ErrRecordNotFound
	}
	return u, nil
}

func (r *SessionRepository) FindByUserID(userID string) (*model.Session, error) {
	for _, s := range r.sessions {
		if s.UserID == userID {
			return s, nil
		}
	}
	return nil, store.ErrRecordNotFound
}

func (r *SessionRepository) Delete(userID string) error {
	for k, s := range r.sessions {
		if s.UserID == userID {
			delete(r.sessions, k)
			return nil
		}
	}
	return store.ErrRecordNotFound
}
