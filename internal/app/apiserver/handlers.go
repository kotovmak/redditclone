package apiserver

import (
	"context"
	"encoding/json"
	"net/http"
	"redditclone/internal/app/model"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func (s *server) setRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New().String()
		w.Header().Set("X-Request-ID", id)
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyRequestID, id)))
	})
}

func (s *server) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := s.logger.WithFields(logrus.Fields{
			"remore_addr": r.RemoteAddr,
			"uri":         r.URL,
			"request_id":  r.Context().Value(ctxKeyRequestID),
			"method":      r.Method,
		})

		logger.Infof("started %s %s", r.Method, r.RequestURI)

		start := time.Now()
		rw := &responseWriter{w, http.StatusOK}
		next.ServeHTTP(rw, r)

		logger.Infof("completed with %d %s in %v", rw.code, http.StatusText(rw.code), time.Now().Sub(start))
	})
}

func (s *server) authenticateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if header == "" {
			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		headerParts := strings.Split(header, " ")
		token := headerParts[1]
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		if len(token) == 0 {
			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		sess, err := s.store.Session().FindByToken(token)
		if sess == nil || err != nil {
			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		jwtToken, err := s.tokenManager.Parse(token)
		if err != nil {
			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}
		u, err := s.store.User().Find(sess.UserID)
		if err != nil {
			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}
		if !jwtToken.Valid {
			s.store.Session().Delete(u.ID)
			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyUser, u)))

	})
}

func (s *server) handleRegister() http.HandlerFunc {
	type Request struct {
		Username string `json:"username" validate:"required,alphanum,max=32"`
		Password string `json:"password" validete:"required,min=8,max=72"`
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
		u1, err := s.store.User().FindByUsername(req.Username)
		if u1 != nil {
			el := []*errList{}
			el = s.makeEL(el, "body", "already exist", "username", req.Username)
			s.respond(w, r, http.StatusUnprocessableEntity, el)
			return
		}

		u := &model.User{
			Username: req.Username,
			Password: req.Password,
		}
		if err := s.store.User().Create(u); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}
		u.Sanitize()

		token, err := s.tokenManager.NewJWT(u)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
		}
		session := &model.Session{
			UserID:    u.ID,
			RequestID: r.Context().Value(ctxKeyRequestID).(string),
			Token:     token,
		}

		if err := s.store.Session().Create(session); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		res := &AuthResponse{
			"token": token,
		}
		s.respond(w, r, http.StatusCreated, res)
	}
}

func (s *server) handleLogin() http.HandlerFunc {
	type Request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &Request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		u, err := s.store.User().FindByUsername(req.Username)
		if err != nil || !u.ComparePassword(req.Password) {
			s.message(w, r, http.StatusUnauthorized, errIncorrectUsernameOrPassword)
			return
		}

		s1, err := s.store.Session().FindByUserID(u.ID)
		if s1 != nil {
			s.store.Session().Delete(s1.UserID)
		}

		token, err := s.tokenManager.NewJWT(u)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
		}
		session := &model.Session{
			UserID:    u.ID,
			RequestID: r.Context().Value(ctxKeyRequestID).(string),
			Token:     token,
		}

		if err := s.store.Session().Create(session); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		res := &AuthResponse{
			"token": token,
		}
		s.respond(w, r, http.StatusOK, res)
	}
}
