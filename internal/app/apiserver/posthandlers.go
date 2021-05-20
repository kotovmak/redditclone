package apiserver

import (
	"errors"
	"net/http"
	"redditclone/internal/app/model"

	"github.com/gorilla/mux"
)

func (s *server) handlePostDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, ok := vars["post_id"]
		if !ok {
			s.error(w, r, http.StatusBadRequest, mux.ErrNotFound)
			return
		}
		p, err := s.store.Post().FindByID(id)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		e1 := s.store.Post().Delete(p.ID)

		if e1 != nil {
			s.error(w, r, http.StatusInternalServerError, e1)
			return
		}
		s.store.Comment().DeleteByPostID(id)
		s.store.Vote().DeleteByPostID(id)
		s.message(w, r, http.StatusOK, errors.New("success"))
	}
}

func (s *server) handleUserPostsList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		pl, err := s.store.Post().UserPostsList(vars["user_login"])
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		comments, err := s.store.Comment().List()
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		votes, err := s.store.Vote().List()
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		respond := []*model.Post{}
		for _, p := range pl {
			v, ok := votes[p.ID]
			if ok {
				p.Votes = v
			} else {
				p.Votes = []*model.Vote{}
			}
			c, ok := comments[p.ID]
			if ok {
				p.Comments = c
			} else {
				p.Comments = []*model.Comment{}
			}
			p.Comments = comments[p.ID]
			respond = append(respond, p)
		}
		s.respond(w, r, http.StatusOK, respond)

	}
}

func (s *server) handleCategoryList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		pl, err := s.store.Post().CategoryList(vars["category_name"])
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		comments, err := s.store.Comment().List()
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		votes, err := s.store.Vote().List()
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		respond := []*model.Post{}
		for _, p := range pl {
			v, ok := votes[p.ID]
			if ok {
				p.Votes = v
			} else {
				p.Votes = []*model.Vote{}
			}
			c, ok := comments[p.ID]
			if ok {
				p.Comments = c
			} else {
				p.Comments = []*model.Comment{}
			}
			p.Comments = comments[p.ID]
			respond = append(respond, p)
		}
		s.respond(w, r, http.StatusOK, respond)

	}
}

func (s *server) handlePostList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pl, err := s.store.Post().List()
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		comments, err := s.store.Comment().List()
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		votes, err := s.store.Vote().List()
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		for _, p := range pl {
			v, ok := votes[p.ID]
			if ok {
				p.Votes = v
			} else {
				p.Votes = []*model.Vote{}
			}
			c, ok := comments[p.ID]
			if ok {
				p.Comments = c
			} else {
				p.Comments = []*model.Comment{}
			}
			p.Comments = comments[p.ID]
		}
		s.respond(w, r, http.StatusOK, pl)

	}
}

func (s *server) handlePostCreate() http.HandlerFunc {
	type Request struct {
		Category string `json:"category" validate:"required"`
		Type     string `json:"type" validate:"required,oneof=link text"`
		Title    string `json:"title" validate:"required"`
		Text     string `json:"text" validate:"required_if=type text,min=4"`
		Url      string `json:"url" validate:"omitempty,required_if=type link,url"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &Request{}
		el, err := s.validator(r, req)
		if el != nil {
			s.respond(w, r, http.StatusUnprocessableEntity, el)
			return
		}
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		u := r.Context().Value(ctxKeyUser).(*model.User)
		v := []*model.Vote{}
		c := []*model.Comment{}
		p := &model.Post{
			Category:         req.Category,
			Type:             req.Type,
			Title:            req.Title,
			Text:             req.Text,
			Score:            0,
			Views:            0,
			Author:           u,
			UpvotePercentage: 0,
			Votes:            v,
			Comments:         c,
		}
		if err := s.store.Post().Create(p); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}
		s.respond(w, r, http.StatusCreated, p)
	}
}

func (s *server) handlePostGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, ok := vars["post_id"]
		if !ok {
			s.error(w, r, http.StatusBadRequest, mux.ErrNotFound)
			return
		}
		p, err := s.store.Post().FindByID(id)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		if p == nil {
			s.message(w, r, http.StatusNotFound, errors.New("invalid post id"))
			return
		}

		comments, err := s.store.Comment().ListByPostID(p.ID)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		p.Comments = comments
		votes, err := s.store.Vote().ListByPostID(p.ID)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		p.Votes = votes
		s.store.Post().Viewed(p)
		s.respond(w, r, http.StatusOK, p)
	}
}
