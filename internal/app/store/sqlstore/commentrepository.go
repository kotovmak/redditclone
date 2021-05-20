package sqlstore

import (
	"database/sql"
	"redditclone/internal/app/model"
	"redditclone/internal/app/store"

	"github.com/google/uuid"
)

type CommentRepository struct {
	store *Store
}

func (r *CommentRepository) Create(c *model.Comment) error {
	return r.store.db.QueryRow(
		`INSERT INTO comments
			(body,author,post_id,id)
		 VALUES 
		 	(
				 $1,
				 (SELECT id FROM users WHERE id=$2),
				 (SELECT id FROM posts WHERE id=$3),
				 $4
			) RETURNING id`,
		c.Body,
		c.Author.ID,
		c.PostID,
		uuid.New().String(),
	).Scan(&c.ID)
}

func (r *CommentRepository) Delete(commentID string) error {
	if err := r.store.db.QueryRow(
		"DELETE FROM comments WHERE id = $1",
		commentID,
	).Err(); err != nil {
		return err
	}
	return nil
}

func (r *CommentRepository) DeleteByPostID(postID string) error {
	if err := r.store.db.QueryRow(
		"DELETE FROM comments WHERE post_id = $1",
		postID,
	).Err(); err != nil {
		return err
	}
	return nil
}

func (r *CommentRepository) List() (map[string][]*model.Comment, error) {
	mp := map[string][]*model.Comment{}
	query :=
		`SELECT 
			comments.id,
			comments.body,
			comments.created_at,
			comments.author,
			users.username,
			comments.post_id
		FROM 
			comments
		LEFT JOIN users ON comments.author=users.id`
	data, err := r.store.db.Query(query)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}
	for data.Next() {
		c := &model.Comment{}
		u := &model.User{}
		err = data.Scan(
			&c.ID,
			&c.Body,
			&c.Created,
			&u.ID,
			&u.Username,
			&c.PostID,
		)
		if err != nil {
			return nil, err
		}
		c.Author = u
		mp[c.PostID] = append(mp[c.PostID], c)
	}
	return mp, nil
}

func (r *CommentRepository) ListByPostID(postID string) ([]*model.Comment, error) {
	mp := []*model.Comment{}
	query :=
		`SELECT 
			comments.id,
			comments.body,
			comments.created_at,
			comments.author,
			users.username
		FROM 
			comments
		LEFT JOIN users ON comments.author=users.id
		WHERE 
			comments.post_id=$1`
	data, err := r.store.db.Query(query, postID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}
	for data.Next() {
		c := &model.Comment{}
		u := &model.User{}
		err = data.Scan(
			&c.ID,
			&c.Body,
			&c.Created,
			&u.ID,
			&u.Username,
		)
		if err != nil {
			return nil, err
		}
		c.Author = u
		mp = append(mp, c)
	}
	return mp, nil
}
