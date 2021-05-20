package teststore_test

import (
	"redditclone/internal/app/model"
	"redditclone/internal/app/store/teststore"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommentsRepository_Create(t *testing.T) {
	s := teststore.New()
	c := model.TestComment(t)
	assert.NoError(t, s.Comment().Create(c))
	assert.NotNil(t, c)
}

func TestCommentsRepository_DeleteByPostID(t *testing.T) {
	s := teststore.New()
	u := model.TestUser(t)
	assert.NoError(t, s.User().Create(u))
	p := model.TestPost(t)
	p.Author = u
	assert.NoError(t, s.Post().Create(p))
	c := model.TestComment(t)
	c.PostID = p.ID
	c.Author = u
	assert.NoError(t, s.Comment().Create(c))

	err := s.Comment().DeleteByPostID(p.ID)
	assert.NoError(t, err)
	mp, err := s.Comment().List()
	pl := mp[p.ID]
	assert.NotContains(t, pl, c)
	assert.NoError(t, err)
}
