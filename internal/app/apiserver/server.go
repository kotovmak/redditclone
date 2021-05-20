package apiserver

import (
	"encoding/json"
	"errors"
	"net/http"
	"redditclone/internal/app/store"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

const (
	ctxKeyUser ctxKey = iota
	ctxKeyRequestID
)

var (
	errIncorrectUsernameOrPassword = errors.New("incorrect username or password")
	errNotAuthenticated            = errors.New("unauthorized")
)

type ctxKey int8

type server struct {
	router       *mux.Router
	logger       *logrus.Logger
	store        store.Store
	tokenManager *TokenManager
	validate     *validator.Validate
}

type errList struct {
	Location string `json:"location"`
	Message  string `json:"msg"`
	Param    string `json:"param"`
	Value    string `json:"value" validation:"omitempty"`
}

func NewServer(store store.Store, tokenManager *TokenManager) *server {
	s := &server{
		router:       mux.NewRouter(),
		logger:       logrus.New(),
		store:        store,
		tokenManager: tokenManager,
		validate:     validator.New(),
	}
	s.configureRouter()
	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) configureRouter() {
	s.router.Use(s.setRequestID)
	s.router.Use(s.logRequest)
	s.router.Use(handlers.CORS(handlers.AllowedOrigins([]string{"*"})))
	s.router.HandleFunc("/api/register", s.handleRegister()).Methods("POST")
	s.router.HandleFunc("/api/login", s.handleLogin()).Methods("POST")
	s.router.HandleFunc("/api/posts/", s.handlePostList()).Methods("GET")
	s.router.HandleFunc("/api/posts/{category_name}", s.handleCategoryList()).Methods("GET")
	s.router.HandleFunc("/api/post/{post_id}", s.handlePostGet()).Methods("GET")
	s.router.HandleFunc("/api/user/{user_login}", s.handleUserPostsList()).Methods("GET")

	private := s.router.PathPrefix("/api").Subrouter()
	private.Use(s.authenticateUser)
	private.HandleFunc("/posts", s.handlePostCreate()).Methods("POST")
	private.HandleFunc("/post/{post_id}", s.handleCommentCreate()).Methods("POST")
	private.HandleFunc("/post/{post_id}", s.handlePostDelete()).Methods("DELETE")
	private.HandleFunc("/post/{post_id}/upvote", s.handlePostUpvote()).Methods("GET")
	private.HandleFunc("/post/{post_id}/unvote", s.handlePostUnvote()).Methods("GET")
	private.HandleFunc("/post/{post_id}/downvote", s.handlePostDownvote()).Methods("GET")
	private.HandleFunc("/post/{post_id}/{comment_id}", s.handleCommentDelete()).Methods("DELETE")

	s.router.PathPrefix("/static").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./web/static/"))))

	s.router.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./web/index.html")
	})
}

func (s *server) validator(r *http.Request, req interface{}) ([]*errList, error) {
	el := []*errList{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return nil, err
	}
	err := s.validate.Struct(req)
	if err == nil {
		return nil, nil
	}
	for _, err := range err.(validator.ValidationErrors) {
		el = s.makeEL(el, "body", err.Error(), strings.ToLower(err.StructField()), err.Value().(string))
	}

	return el, nil
}

func (s *server) makeEL(el []*errList, l string, m string, p string, v string) []*errList {
	e := &errList{}
	e.Location = l
	e.Message = m
	e.Param = p
	e.Value = v
	el = append(el, e)
	return el
}

func (s *server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{"error": err.Error()})
}

func (s *server) message(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{"message": err.Error()})
}

func (s *server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
