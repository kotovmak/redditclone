package teststore

import (
	"redditclone/internal/app/model"
	"redditclone/internal/app/store"
)

type Store struct {
	userRepository    *UserRepository
	sessionRepository *SessionRepository
	postRepository    *PostRepository
	commentRepository *CommentRepository
	voteRepository    *VoteRepository
}

func New() *Store {
	return &Store{}
}

func (s *Store) User() store.UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}

	s.userRepository = &UserRepository{
		store: s,
		users: make(map[string]*model.User),
	}

	return s.userRepository
}

func (s *Store) Session() store.SessionRepository {
	if s.sessionRepository != nil {
		return s.sessionRepository
	}

	s.sessionRepository = &SessionRepository{
		store:    s,
		sessions: make(map[string]*model.Session),
	}

	return s.sessionRepository
}

func (s *Store) Post() store.PostRepository {
	if s.postRepository != nil {
		return s.postRepository
	}

	s.postRepository = &PostRepository{
		store: s,
		posts: make(map[string]*model.Post),
	}

	return s.postRepository
}

func (s *Store) Comment() store.CommentRepository {
	if s.commentRepository != nil {
		return s.commentRepository
	}

	s.commentRepository = &CommentRepository{
		store:    s,
		comments: make(map[string]*model.Comment),
	}

	return s.commentRepository
}

func (s *Store) Vote() store.VoteRepository {
	if s.voteRepository != nil {
		return s.voteRepository
	}

	s.voteRepository = &VoteRepository{
		store: s,
		votes: make(map[string]*model.Vote),
	}

	return s.voteRepository
}
