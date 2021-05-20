package sqlstore

import (
	"database/sql"
	"redditclone/internal/app/model"
	"redditclone/internal/app/store"

	"github.com/google/uuid"
)

type PostRepository struct {
	store *Store
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
	return r.store.db.QueryRow(
		`UPDATE posts SET
			(score, upvote_percentage) =
		 	(
				 $1,
				 $2
			)
		WHERE
			id=$3
			`,
		p.Score,
		p.UpvotePercentage,
		p.ID,
	).Err()
}

func (r *PostRepository) Viewed(p *model.Post) error {
	p.Views++
	return r.store.db.QueryRow(
		`UPDATE posts SET
			views = $1
		WHERE
			id=$2
			`,
		p.Views,
		p.ID,
	).Err()
}

func (r *PostRepository) Create(p *model.Post) error {
	return r.store.db.QueryRow(
		`INSERT INTO posts 
			(score,views,type,title,author,category,text,id,upvote_percentage)
		 VALUES 
		 	(
				 $1,
				 $2,
				 $3,
				 $4,
				 (SELECT id FROM users WHERE id=$5),
				 $6,
				 $7,
				 $8,
				 $9
			) RETURNING id`,
		p.Score,
		p.Views,
		p.Type,
		p.Title,
		p.Author.ID,
		p.Category,
		p.Text,
		uuid.New().String(),
		p.UpvotePercentage,
	).Scan(&p.ID)
}

func (r *PostRepository) List() ([]*model.Post, error) {
	mp := []*model.Post{}
	query :=
		`SELECT 
			posts.id,
			posts.score,
			posts.views,
			posts.type,
			posts.title,
			posts.author,
			users.username,
			posts.category,
			posts.text,
			posts.created_at 
		FROM 
			posts
		LEFT JOIN users ON posts.author=users.id`
	data, err := r.store.db.Query(query)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}
	for data.Next() {
		p := &model.Post{}
		u := &model.User{}
		err = data.Scan(
			&p.ID,
			&p.Score,
			&p.Views,
			&p.Type,
			&p.Title,
			&u.ID,
			&u.Username,
			&p.Category,
			&p.Text,
			&p.Created,
		)
		if err != nil {
			return nil, err
		}
		p.Author = u
		mp = append(mp, p)
	}
	return mp, nil
}

func (r *PostRepository) CategoryList(categoryName string) ([]*model.Post, error) {
	mp := []*model.Post{}
	query :=
		`SELECT 
			posts.id,
			posts.score,
			posts.views,
			posts.type,
			posts.title,
			posts.author,
			users.username,
			posts.category,
			posts.text,
			posts.created_at 
		FROM 
			posts
		LEFT JOIN users ON posts.author=users.id
		WHERE posts.category=$1`
	data, err := r.store.db.Query(query, categoryName)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}
	for data.Next() {
		p := &model.Post{}
		u := &model.User{}
		err = data.Scan(
			&p.ID,
			&p.Score,
			&p.Views,
			&p.Type,
			&p.Title,
			&u.ID,
			&u.Username,
			&p.Category,
			&p.Text,
			&p.Created,
		)
		if err != nil {
			return nil, err
		}
		p.Author = u
		mp = append(mp, p)
	}
	return mp, nil
}

func (r *PostRepository) UserPostsList(username string) ([]*model.Post, error) {
	mp := []*model.Post{}
	query :=
		`SELECT 
			posts.id,
			posts.score,
			posts.views,
			posts.type,
			posts.title,
			posts.category,
			posts.text,
			posts.created_at 
		FROM 
			posts
		WHERE author=(SELECT id FROM users WHERE username=$1)`
	data, err := r.store.db.Query(query, username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}
	for data.Next() {
		p := &model.Post{}
		err = data.Scan(
			&p.ID,
			&p.Score,
			&p.Views,
			&p.Type,
			&p.Title,
			&p.Category,
			&p.Text,
			&p.Created,
		)
		if err != nil {
			return nil, err
		}
		mp = append(mp, p)
	}
	return mp, nil
}

func (r *PostRepository) FindByID(postID string) (*model.Post, error) {
	p := &model.Post{}
	u := &model.User{}
	if err := r.store.db.QueryRow(
		`SELECT 
			posts.id,
			posts.score,
			posts.views,
			posts.type,
			posts.title,
			posts.author,
			users.username,
			posts.category,
			posts.text,
			posts.created_at 
		FROM 
			posts
		LEFT JOIN users ON posts.author=users.id
		WHERE 
			posts.id = $1`,
		postID,
	).Scan(
		&p.ID,
		&p.Score,
		&p.Views,
		&p.Type,
		&p.Title,
		&u.ID,
		&u.Username,
		&p.Category,
		&p.Text,
		&p.Created,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}
	p.Author = u
	comments, err := r.store.Comment().ListByPostID(p.ID)
	if err != nil {
		return nil, err
	}
	p.Comments = comments
	votes, err := r.store.Vote().ListByPostID(p.ID)
	if err != nil {
		return nil, err
	}
	p.Votes = votes
	return p, nil
}

func (r *PostRepository) Delete(postID string) error {
	if err := r.store.db.QueryRow(
		"DELETE FROM posts WHERE id = $1",
		postID,
	).Err(); err != nil {
		return err
	}
	return nil
}
