package model

type Vote struct {
	ID     string `json:"-"`
	UserID string `json:"user"`
	Vote   int    `json:"vote"`
	PostID string `json:"-"`
}
