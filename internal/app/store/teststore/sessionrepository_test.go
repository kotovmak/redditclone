package teststore_test

import (
	"redditclone/internal/app/model"
	"redditclone/internal/app/store/teststore"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSessionRepository_Create(t *testing.T) {
	s := teststore.New()
	sess := model.TestSession(t)
	assert.NoError(t, s.Session().Create(sess))
	assert.NotNil(t, sess)
}

func TestSessionRepository_FindByToken(t *testing.T) {
	s := teststore.New()
	s1 := model.TestSession(t)
	u := model.TestUser(t)
	s.User().Create(u)
	s1.UserID = u.ID
	s.Session().Create(s1)
	u2, err := s.Session().FindByToken(s1.Token)
	assert.NoError(t, err)
	assert.NotNil(t, u2)
}

func TestSessionRepository_FindByUserID(t *testing.T) {
	s := teststore.New()
	s1 := model.TestSession(t)
	s.Session().Create(s1)
	u2, err := s.Session().FindByUserID(s1.UserID)
	assert.NoError(t, err)
	assert.NotNil(t, u2)
}
