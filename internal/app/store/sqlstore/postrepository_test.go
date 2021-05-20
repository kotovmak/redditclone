package sqlstore_test

import (
	"redditclone/internal/app/model"
	"redditclone/internal/app/store/sqlstore"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPostRepository_Create(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("posts", "users")
	s := sqlstore.New(db)
	u := model.TestUser(t)
	assert.NoError(t, s.User().Create(u))
	assert.NotNil(t, u)
	p := model.TestPost(t)
	p.Author = u
	assert.NoError(t, s.Post().Create(p))
	assert.NotNil(t, p)
}

func TestPostRepository_FindByID(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("posts", "users")

	s := sqlstore.New(db)
	u := model.TestUser(t)
	s.User().Create(u)
	p := model.TestPost(t)
	p.Author = u
	s.Post().Create(p)
	p2, err := s.Post().FindByID(p.ID)
	assert.NoError(t, err)
	assert.NotNil(t, p2)
}

func TestPostRepository_Recalc(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("posts", "users", "votes")

	s := sqlstore.New(db)
	u := model.TestUser(t)
	assert.NoError(t, s.User().Create(u))
	assert.NotNil(t, u)

	p := model.TestPost(t)
	p.Author = u
	assert.NoError(t, s.Post().Create(p))
	assert.NotNil(t, p)

	v := model.TestVote(t)
	v.PostID = p.ID
	v.UserID = u.ID
	assert.NoError(t, s.Vote().Create(v))
	assert.NotNil(t, p)

	p.Votes = append(p.Votes, v)
	s.Post().Recalc(p)

	p2, err := s.Post().FindByID(p.ID)
	assert.NoError(t, err)
	assert.NotNil(t, p2)
	assert.Equal(t, v.Vote, p2.Score)
}
