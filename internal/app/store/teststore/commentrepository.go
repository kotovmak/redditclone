package teststore

import (
	"redditclone/internal/app/model"
	"redditclone/internal/app/store"

	"github.com/google/uuid"
)

type CommentRepository struct {
	store    *Store
	comments map[string]*model.Comment
}

func (r *CommentRepository) Create(c *model.Comment) error {
	c.ID = uuid.New().String()
	r.comments[c.ID] = c
	return nil
}

func (r *CommentRepository) Delete(commentID string) error {
	_, ok := r.comments[commentID]
	if !ok {
		return store.ErrRecordNotFound
	}
	delete(r.comments, commentID)
	return nil
}

func (r *CommentRepository) DeleteByPostID(postID string) error {
	for _, s := range r.comments {
		if s.PostID == postID {
			delete(r.comments, s.ID)
			return nil
		}
	}
	return store.ErrRecordNotFound
}

func (r *CommentRepository) List() (map[string][]*model.Comment, error) {
	mp := map[string][]*model.Comment{}
	for _, s := range r.comments {
		mp[s.PostID] = append(mp[s.PostID], s)
	}
	return mp, nil
}

func (r *CommentRepository) ListByPostID(postID string) ([]*model.Comment, error) {
	mp := []*model.Comment{}
	for _, s := range r.comments {
		if s.PostID == postID {
			mp = append(mp, s)
		}
	}
	return mp, nil
}
