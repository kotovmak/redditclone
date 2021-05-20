package sqlstore

import (
	"database/sql"
	"redditclone/internal/app/model"
	"redditclone/internal/app/store"

	"github.com/google/uuid"
)

type VoteRepository struct {
	store *Store
}

func (r *VoteRepository) Create(v *model.Vote) error {
	return r.store.db.QueryRow(
		`INSERT INTO votes
			("vote","post_id","user","id")
		 VALUES 
		 	(
				 $1,
				 (SELECT id FROM posts WHERE id=$2),
				 (SELECT id FROM users WHERE id=$3),
				 $4
			) RETURNING id`,
		v.Vote,
		v.PostID,
		v.UserID,
		uuid.New().String(),
	).Scan(&v.ID)
}

func (r *VoteRepository) Delete(voteID string) error {
	if err := r.store.db.QueryRow(
		"DELETE FROM votes WHERE id = $1",
		voteID,
	).Err(); err != nil {
		return err
	}
	return nil
}

func (r *VoteRepository) DeleteByPostID(postID string) error {
	if err := r.store.db.QueryRow(
		"DELETE FROM votes WHERE post_id = $1",
		postID,
	).Err(); err != nil {
		return err
	}
	return nil
}

func (r *VoteRepository) List() (map[string][]*model.Vote, error) {
	mp := map[string][]*model.Vote{}
	query :=
		`SELECT 
			votes.id,
			votes.vote,
			votes."user",
			votes.post_id
		FROM 
			votes`
	data, err := r.store.db.Query(query)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}
	for data.Next() {
		v := &model.Vote{}
		err = data.Scan(
			&v.ID,
			&v.Vote,
			&v.UserID,
			&v.PostID,
		)
		if err != nil {
			return nil, err
		}
		mp[v.PostID] = append(mp[v.PostID], v)
	}
	return mp, nil
}

func (r *VoteRepository) ListByPostID(postID string) ([]*model.Vote, error) {
	mp := []*model.Vote{}
	query :=
		`SELECT 
			votes.id,
			votes.vote,
			votes."user",
			votes.post_id
		FROM 
			votes
		WHERE 
			votes.post_id=$1`
	data, err := r.store.db.Query(query, postID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}
	for data.Next() {
		v := &model.Vote{}
		err = data.Scan(
			&v.ID,
			&v.Vote,
			&v.UserID,
			&v.PostID,
		)
		if err != nil {
			return nil, err
		}
		mp = append(mp, v)
	}
	return mp, nil
}
