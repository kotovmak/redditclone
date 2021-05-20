package apiserver

import (
	"net/http"
	"redditclone/internal/app/model"

	"github.com/gorilla/mux"
)

func (s *server) handleCommentDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		commentID, ok := vars["comment_id"]
		if !ok {
			s.error(w, r, http.StatusBadRequest, mux.ErrNotFound)
			return
		}
		postID, ok := vars["post_id"]
		if !ok {
			s.error(w, r, http.StatusBadRequest, mux.ErrNotFound)
			return
		}
		c := s.store.Comment().Delete(commentID)
		if c != nil {
			s.error(w, r, http.StatusInternalServerError, c)
			return
		}
		p, err := s.store.Post().FindByID(postID)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, c)
			return
		}
		s.respond(w, r, http.StatusOK, p)
	}
}

func (s *server) handleCommentCreate() http.HandlerFunc {
	type Request struct {
		Comment string `json:"comment" validate:"required"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &Request{}
		el, err := s.validator(r, req)
		if el != nil {
			s.respond(w, r, http.StatusUnprocessableEntity, el)
			return
		}
		if err != nil {
			s.error(w, r, http.StatusBadRequest, http.ErrBodyNotAllowed)
			return
		}
		vars := mux.Vars(r)
		post_id, ok := vars["post_id"]
		if !ok {
			s.error(w, r, http.StatusBadRequest, mux.ErrNotFound)
			return
		}
		u := r.Context().Value(ctxKeyUser).(*model.User)
		c := &model.Comment{
			Body:   req.Comment,
			PostID: post_id,
			Author: u,
		}
		if err := s.store.Comment().Create(c); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}
		p, err := s.store.Post().FindByID(post_id)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusCreated, p)
	}
}
