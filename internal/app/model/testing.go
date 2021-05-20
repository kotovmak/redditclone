package model

import "testing"

func TestUser(t *testing.T) *User {
	return &User{
		Username: "user",
		Password: "password",
	}
}

func TestSession(t *testing.T) *Session {
	return &Session{
		RequestID: "a4b93a56-034a-4044-8cd1-1acf08bc5127",
		Token:     "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjp7InVzZXJuYW1lIjoicXdlNTY3NTY3NTY3IiwiaWQiOiI2MGEyMzU1MTNiNzQ0NjAwMDg2MjQxYzIifSwiaWF0IjoxNjIxMjQzMjE3LCJleHAiOjE2MjE4NDgwMTd9.XfXawaGg_jygYHiGh42dZSkNBge1SDmazOChtnN5chA",
	}
}

func TestPost(t *testing.T) *Post {
	return &Post{
		Category: "music",
		Type:     "text",
		Title:    "Header",
		Text:     "some text",
	}
}

func TestComment(t *testing.T) *Comment {
	return &Comment{
		Body: "test comment",
	}
}

func TestVote(t *testing.T) *Vote {
	return &Vote{
		Vote: 1,
	}
}
