package store

import "redditclone/internal/app/model"

type UserRepository interface {
	Create(*model.User) error
	Find(string) (*model.User, error)
	FindByUsername(string) (*model.User, error)
}

type SessionRepository interface {
	Create(*model.Session) error
	FindByToken(string) (*model.Session, error)
	FindByUserID(string) (*model.Session, error)
	Delete(string) error
}

type PostRepository interface {
	Create(*model.Post) error
	Recalc(*model.Post) error
	Viewed(*model.Post) error
	List() ([]*model.Post, error)
	CategoryList(string) ([]*model.Post, error)
	UserPostsList(string) ([]*model.Post, error)
	FindByID(string) (*model.Post, error)
	Delete(string) error
}

type CommentRepository interface {
	List() (map[string][]*model.Comment, error)
	ListByPostID(string) ([]*model.Comment, error)
	Create(*model.Comment) error
	Delete(string) error
	DeleteByPostID(string) error
}

type VoteRepository interface {
	List() (map[string][]*model.Vote, error)
	ListByPostID(string) ([]*model.Vote, error)
	Create(*model.Vote) error
	Delete(string) error
	DeleteByPostID(string) error
}
