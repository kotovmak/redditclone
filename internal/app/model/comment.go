package model

import "time"

type Comment struct {
	Author  *User     `json:"author"`
	Body    string    `json:"body"`
	Created time.Time `json:"created"`
	ID      string    `json:"id"`
	PostID  string    `json:"-"`
}
