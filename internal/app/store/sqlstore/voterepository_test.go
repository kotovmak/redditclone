package sqlstore_test

import (
	"redditclone/internal/app/model"
	"redditclone/internal/app/store/sqlstore"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVotesRepository_Create(t *testing.T) {
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
	v.UserID = u.ID
	v.PostID = p.ID
	assert.NoError(t, s.Vote().Create(v))
	assert.NotNil(t, p)
}
