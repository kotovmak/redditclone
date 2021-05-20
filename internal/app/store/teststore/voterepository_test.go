package teststore_test

import (
	"redditclone/internal/app/model"
	"redditclone/internal/app/store/teststore"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVoteRepository_Create(t *testing.T) {
	s := teststore.New()
	v := model.TestVote(t)
	assert.NoError(t, s.Vote().Create(v))
	assert.NotNil(t, v)
}

func TestVoteRepository_DeleteByPostID(t *testing.T) {
	s := teststore.New()
	u := model.TestUser(t)
	assert.NoError(t, s.User().Create(u))
	p := model.TestPost(t)
	p.Author = u
	assert.NoError(t, s.Post().Create(p))
	v := model.TestVote(t)
	v.PostID = p.ID
	v.UserID = u.ID
	assert.NoError(t, s.Vote().Create(v))

	err := s.Vote().DeleteByPostID(p.ID)
	assert.NoError(t, err)
	mp, err := s.Vote().List()
	pl := mp[p.ID]
	assert.NotContains(t, pl, v)
	assert.NoError(t, err)
}
