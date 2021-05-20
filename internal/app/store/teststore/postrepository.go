package teststore

import (
	"redditclone/internal/app/model"
	"redditclone/internal/app/store"

	"github.com/google/uuid"
)

type PostRepository struct {
	store *Store
	posts map[string]*model.Post
}

func (r *PostRepository) Recalc(p *model.Post) error {
	var up, down int
	p.Score = 0
	for _, c := range p.Votes {
		if c.Vote > 0 {
			up++
		} else {
			down++
		}
		p.Score += c.Vote
	}
	if up > 0 && up > down {
		p.UpvotePercentage = len(p.Votes) / up * 100
	}
	r.posts[p.ID].Score = p.Score
	r.posts[p.ID].UpvotePercentage = p.UpvotePercentage
	return nil
}

func (r *PostRepository) Viewed(p *model.Post) error {
	r.posts[p.ID].Views++
	return nil
}

func (r *PostRepository) Create(p *model.Post) error {
	p.ID = uuid.New().String()
	r.posts[p.ID] = p
	return nil
}

func (r *PostRepository) List() ([]*model.Post, error) {
	mp := []*model.Post{}
	for _, p := range r.posts {
		mp = append(mp, p)
	}
	return mp, nil
}

func (r *PostRepository) FindByID(postID string) (*model.Post, error) {
	u, ok := r.posts[postID]
	if !ok {
		return nil, store.ErrRecordNotFound
	}
	c, _ := r.store.Comment().ListByPostID(postID)
	r.posts[postID].Comments = c
	v, _ := r.store.Vote().ListByPostID(postID)
	r.posts[postID].Votes = v
	return u, nil
}

func (r *PostRepository) Delete(postID string) error {
	for k, p := range r.posts {
		if p.ID == postID {
			delete(r.posts, k)
		}
	}
	return nil
}

func (r *PostRepository) CategoryList(categoryName string) ([]*model.Post, error) {
	mp := []*model.Post{}
	for _, p := range r.posts {
		if p.Category == categoryName {
			mp = append(mp, p)
		}
	}
	return mp, nil
}

func (r *PostRepository) UserPostsList(username string) ([]*model.Post, error) {
	mp := []*model.Post{}
	for _, p := range r.posts {
		if p.Author.Username == username {
			mp = append(mp, p)
		}
	}
	return mp, nil
}
