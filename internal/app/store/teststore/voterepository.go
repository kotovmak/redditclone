package teststore

import (
	"redditclone/internal/app/model"
	"redditclone/internal/app/store"

	"github.com/google/uuid"
)

type VoteRepository struct {
	store *Store
	votes map[string]*model.Vote
}

func (r *VoteRepository) Create(v *model.Vote) error {
	v.ID = uuid.New().String()
	r.votes[v.ID] = v
	return nil
}

func (r *VoteRepository) Delete(voteID string) error {
	_, ok := r.votes[voteID]
	if !ok {
		return store.ErrRecordNotFound
	}
	delete(r.votes, voteID)
	return nil
}

func (r *VoteRepository) DeleteByPostID(postID string) error {
	for _, s := range r.votes {
		if s.PostID == postID {
			delete(r.votes, s.ID)
			return nil
		}
	}
	return store.ErrRecordNotFound
}

func (r *VoteRepository) List() (map[string][]*model.Vote, error) {
	mp := map[string][]*model.Vote{}
	for _, s := range r.votes {
		mp[s.PostID] = append(mp[s.PostID], s)
	}
	return mp, nil
}

func (r *VoteRepository) ListByPostID(postID string) ([]*model.Vote, error) {
	mp := []*model.Vote{}
	for _, s := range r.votes {
		if s.PostID == postID {
			mp = append(mp, s)
		}
	}
	return mp, nil
}
