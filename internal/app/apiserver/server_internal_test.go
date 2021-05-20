package apiserver

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"redditclone/internal/app/model"
	"redditclone/internal/app/store/teststore"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServer_HandleRegister(t *testing.T) {
	tm, _ := NewManager("test_secret_string", "1h")

	s := NewServer(teststore.New(), tm)
	testCases := []struct {
		name         string
		payload      interface{}
		expectedCode int
	}{
		{
			name: "valid",
			payload: map[string]string{
				"username": "user",
				"password": "password",
			},
			expectedCode: http.StatusCreated,
		},
		{
			name:         "invalid payload",
			payload:      "invalid",
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "invalid params",
			payload: map[string]string{
				"username": "u",
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			b := &bytes.Buffer{}
			json.NewEncoder(b).Encode(tc.payload)
			req, _ := http.NewRequest(http.MethodPost, "/api/register", b)
			s.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)

		})
	}
}

func TestServer_HandlePostCreate(t *testing.T) {
	store := teststore.New()
	tm, _ := NewManager("test_secret_string", "1h")
	s := NewServer(store, tm)
	u := model.TestUser(t)
	store.User().Create(u)
	p := model.TestPost(t)
	p.Author = u
	store.Post().Create(p)
	token, _ := tm.NewJWT(u)
	session := &model.Session{
		UserID: u.ID,
		Token:  token,
	}
	s.store.Session().Create(session)

	testCases := []struct {
		name         string
		payload      interface{}
		expectedCode int
	}{
		{
			name: "valid",
			payload: map[string]string{
				"type":     p.Type,
				"category": p.Category,
				"title":    p.Title,
				"text":     p.Text,
			},
			expectedCode: http.StatusCreated,
		},
		{
			name:         "invalid payload",
			payload:      "invalid",
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "text empty",
			payload: map[string]string{
				"type":     p.Type,
				"category": p.Category,
				"title":    p.Title,
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "type empty",
			payload: map[string]string{
				"text":     p.Text,
				"category": p.Category,
				"title":    p.Title,
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "category empty",
			payload: map[string]string{
				"type":  p.Type,
				"text":  p.Text,
				"title": p.Title,
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "title empty",
			payload: map[string]string{
				"type":     p.Type,
				"category": p.Category,
				"text":     p.Text,
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			b := &bytes.Buffer{}
			json.NewEncoder(b).Encode(tc.payload)
			req, _ := http.NewRequest(http.MethodPost, "/api/posts", b)
			req.Header.Set("Authorization", "Bearer "+token)
			s.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)

		})
	}
}

func TestServer_HandlePostList(t *testing.T) {
	store := teststore.New()
	tm, _ := NewManager("test_secret_string", "1h")
	s := NewServer(store, tm)
	u := model.TestUser(t)
	store.User().Create(u)
	p := model.TestPost(t)
	p.Author = u
	store.Post().Create(p)
	u2 := model.TestUser(t)
	u2.Username = "User2"
	store.User().Create(u2)
	p2 := model.TestPost(t)
	p2.Author = u2
	p2.Category = "funny"
	store.Post().Create(p2)
	list, _ := store.Post().List()
	category, _ := store.Post().CategoryList(p.Category)
	user, _ := store.Post().UserPostsList(u.Username)

	testCases := []struct {
		name     string
		payload  []*model.Post
		endpoint string
		count    int
	}{
		{
			name:     "List",
			payload:  list,
			endpoint: "/api/posts/",
			count:    2,
		},
		{
			name:     "Category",
			payload:  category,
			endpoint: "/api/posts/" + p.Category,
			count:    1,
		},
		{
			name:     "User",
			payload:  user,
			endpoint: "/api/user/" + u.Username,
			count:    1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			payload := []*model.Post{}
			b := &bytes.Buffer{}
			json.NewEncoder(b).Encode(tc.payload)
			req, _ := http.NewRequest(http.MethodGet, tc.endpoint, b)
			s.ServeHTTP(rec, req)

			err := json.NewDecoder(rec.Body).Decode(&payload)
			assert.NoError(t, err)
			assert.Equal(t, len(tc.payload), tc.count)
			assert.Equal(t, http.StatusOK, rec.Code)
		})
	}
}

func TestServer_HandlePostGet(t *testing.T) {
	store := teststore.New()
	tm, _ := NewManager("test_secret_string", "1h")
	s := NewServer(store, tm)
	u := model.TestUser(t)
	store.User().Create(u)
	p := model.TestPost(t)
	p.Author = u
	store.Post().Create(p)
	token, _ := tm.NewJWT(u)
	session := &model.Session{
		UserID: u.ID,
		Token:  token,
	}
	s.store.Session().Create(session)

	testCases := []struct {
		name       string
		score      int
		endpoint   string
		countVotes int
	}{
		{
			name:       "Post GET by ID",
			score:      0,
			endpoint:   "/api/post/" + p.ID,
			countVotes: 0,
		},
		{
			name:       "Upvote",
			score:      1,
			endpoint:   "/api/post/" + p.ID + "/upvote",
			countVotes: 1,
		},
		{
			name:       "Unvote",
			score:      0,
			endpoint:   "/api/post/" + p.ID + "/unvote",
			countVotes: 0,
		},
		{
			name:       "Downvote",
			score:      -1,
			endpoint:   "/api/post/" + p.ID + "/downvote",
			countVotes: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			payload := &model.Post{}
			b := &bytes.Buffer{}
			req, _ := http.NewRequest(http.MethodGet, tc.endpoint, b)
			req.Header.Set("Authorization", "Bearer "+token)
			s.ServeHTTP(rec, req)

			err := json.NewDecoder(rec.Body).Decode(&payload)
			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, tc.score, payload.Score)
			assert.Equal(t, tc.countVotes, len(payload.Votes))
		})
	}
}

func TestServer_HandleCommentCreate(t *testing.T) {
	store := teststore.New()
	tm, _ := NewManager("test_secret_string", "1h")
	s := NewServer(store, tm)
	u := model.TestUser(t)
	store.User().Create(u)
	p := model.TestPost(t)
	p.Author = u
	store.Post().Create(p)
	c := model.TestComment(t)
	c.Author = u
	c.PostID = p.ID
	s.store.Comment().Create(c)
	c2 := model.TestComment(t)

	token, _ := tm.NewJWT(u)
	session := &model.Session{
		UserID: u.ID,
		Token:  token,
	}
	s.store.Session().Create(session)

	testCases := []struct {
		name          string
		endpoint      string
		method        string
		payload       interface{}
		countComments int
		code          int
	}{
		{
			name:          "Create comment",
			method:        http.MethodPost,
			endpoint:      "/api/post/" + p.ID,
			countComments: 2,
			payload: map[string]string{
				"comment": c2.Body,
			},
			code: http.StatusCreated,
		},
		{
			name:          "Delete comment",
			method:        http.MethodDelete,
			endpoint:      "/api/post/" + p.ID + "/" + c.ID,
			countComments: 1,
			payload:       "",
			code:          http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			payload := &model.Post{}
			b := &bytes.Buffer{}
			json.NewEncoder(b).Encode(tc.payload)
			req, _ := http.NewRequest(tc.method, tc.endpoint, b)
			req.Header.Set("Authorization", "Bearer "+token)
			s.ServeHTTP(rec, req)

			err := json.NewDecoder(rec.Body).Decode(&payload)
			assert.NoError(t, err)
			assert.Equal(t, tc.code, rec.Code)
			assert.Equal(t, tc.countComments, len(payload.Comments))
		})
	}
}

func TestServer_HandlePostDelete(t *testing.T) {
	store := teststore.New()
	tm, _ := NewManager("test_secret_string", "1h")
	s := NewServer(store, tm)
	u := model.TestUser(t)
	store.User().Create(u)
	p := model.TestPost(t)
	p.Author = u
	store.Post().Create(p)

	token, _ := tm.NewJWT(u)
	session := &model.Session{
		UserID: u.ID,
		Token:  token,
	}
	s.store.Session().Create(session)

	testCases := []struct {
		name     string
		endpoint string
		method   string
		response interface{}
		code     int
	}{
		{
			name:     "Delete post",
			method:   http.MethodDelete,
			endpoint: "/api/post/" + p.ID,
			response: map[string]string{
				"message": "success",
			},
			code: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			response := map[string]string{}
			b := &bytes.Buffer{}
			req, _ := http.NewRequest(tc.method, tc.endpoint, b)
			req.Header.Set("Authorization", "Bearer "+token)
			s.ServeHTTP(rec, req)

			err := json.NewDecoder(rec.Body).Decode(&response)
			assert.NoError(t, err)
			assert.Equal(t, tc.code, rec.Code)
			assert.Equal(t, tc.response, response)
		})
	}
}
