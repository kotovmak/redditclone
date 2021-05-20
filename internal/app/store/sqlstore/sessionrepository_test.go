package sqlstore_test

import (
	"redditclone/internal/app/model"
	"redditclone/internal/app/store"
	"redditclone/internal/app/store/sqlstore"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSessionRepository_Create(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("sessions", "users")
	s := sqlstore.New(db)
	u := model.TestUser(t)
	s.User().Create(u)
	sess := model.TestSession(t)
	sess.UserID = u.ID
	assert.NoError(t, s.Session().Create(sess))
	assert.NotNil(t, sess)
}

func TestSessionRepository_FindByRequestID(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("sessions", "users")

	s := sqlstore.New(db)
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjp7InVzZXJuYW1lIjoicXdlNTY3NTY3NTY3IiwiaWQiOiI2MGEyMzU1MTNiNzQ0NjAwMDg2MjQxYzIifSwiaWF0IjoxNjIxMjQzMjE3LCJleHAiOjE2MjE4NDgwMTd9.XfXawaGg_jygYHiGh42dZSkNBge1SDmazOChtnN5chA"
	_, err := s.Session().FindByToken(token)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	u := model.TestUser(t)
	s.User().Create(u)
	sess := model.TestSession(t)
	sess.UserID = u.ID
	sess.RequestID = token
	s.Session().Create(sess)
	s2, err := s.Session().FindByToken(token)
	assert.NoError(t, err)
	assert.NotNil(t, s2)
}
