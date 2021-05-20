package apiserver

import (
	"net/http"
	"redditclone/internal/app/model"

	"github.com/gorilla/mux"
)

func (s *server) handlePostUpvote() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		post_id, ok := vars["post_id"]
		if !ok {
			s.error(w, r, http.StatusBadRequest, mux.ErrNotFound)
			return
		}
		u := r.Context().Value(ctxKeyUser).(*model.User)

		p, err := s.store.Post().FindByID(post_id)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		for k, v := range p.Votes {
			if v.UserID == u.ID {
				if v.Vote == 1 {
					s.respond(w, r, http.StatusOK, p)
					return
				} else {
					p.Votes = append(p.Votes[:k], p.Votes[k+1:]...)
					s.store.Vote().Delete(v.ID)
				}
			}
		}
		v := &model.Vote{
			Vote:   1,
			PostID: p.ID,
			UserID: u.ID,
		}
		if err := s.store.Vote().Create(v); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}
		p.Votes = append(p.Votes, v)
		s.store.Post().Recalc(p)
		s.respond(w, r, http.StatusOK, p)
	}
}

func (s *server) handlePostDownvote() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, ok := vars["post_id"]
		if !ok {
			s.error(w, r, http.StatusBadRequest, mux.ErrNotFound)
			return
		}
		u := r.Context().Value(ctxKeyUser).(*model.User)

		p, err := s.store.Post().FindByID(id)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		for k, v := range p.Votes {
			if v.UserID == u.ID {
				if v.Vote == -1 {
					s.respond(w, r, http.StatusOK, p)
					return
				} else {
					p.Votes = append(p.Votes[:k], p.Votes[k+1:]...)
					s.store.Vote().Delete(v.ID)
				}
			}
		}
		v := &model.Vote{
			Vote:   -1,
			PostID: p.ID,
			UserID: u.ID,
		}
		if err := s.store.Vote().Create(v); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}
		p.Votes = append(p.Votes, v)
		s.store.Post().Recalc(p)
		s.respond(w, r, http.StatusOK, p)
	}
}

func (s *server) handlePostUnvote() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, ok := vars["post_id"]
		if !ok {
			s.error(w, r, http.StatusBadRequest, mux.ErrNotFound)
			return
		}
		u := r.Context().Value(ctxKeyUser).(*model.User)

		p, err := s.store.Post().FindByID(id)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		for k, v := range p.Votes {
			if v.UserID == u.ID {
				p.Votes = append(p.Votes[:k], p.Votes[k+1:]...)
				s.store.Vote().Delete(v.ID)
				s.store.Post().Recalc(p)
				s.respond(w, r, http.StatusOK, p)
				return
			}
		}
		s.respond(w, r, http.StatusOK, p)
	}
}
