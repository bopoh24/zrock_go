package apiserver

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"reflect"
	"strings"

	"github.com/bopoh24/zrock_go/internal/app/model"
	"github.com/bopoh24/zrock_go/internal/app/settings"
	"github.com/bopoh24/zrock_go/internal/app/store"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// Server api
type Server struct {
	router *mux.Router
	store  store.IfaceStore
	logger *logrus.Logger
}

const (
	contextKeyUserID = "user_id"
)

// NewServer returs server instance
func NewServer(store store.IfaceStore) *Server {
	srv := &Server{
		router: mux.NewRouter(),
		logger: logrus.New(),
		store:  store,
	}
	if level, err := logrus.ParseLevel(settings.App.LogLevel); err == nil {
		srv.logger.SetLevel(level)
	}
	srv.configureRouter()
	srv.logger.Infof("Server started at %s ...", settings.App.BindAdd)
	return srv
}

func (s *Server) configureRouter() {
	//allow CORS
	s.router.Use(handlers.CORS(handlers.AllowedOrigins([]string{"*"})))
	s.router.HandleFunc("/", s.handleHome()).Methods(http.MethodGet)
	authRoutes := s.router.PathPrefix("/api/auth").Subrouter()
	authRoutes.HandleFunc("/login", s.handleLogin()).Methods(http.MethodPost)
	authRoutes.HandleFunc("/register", s.handleRegister()).Methods(http.MethodPost)

	// auth required urls
	api := s.router.PathPrefix("/api").Subrouter()
	api.Use(s.authRequiredJWT)
	api.HandleFunc("/private", s.handlePrivate()).Methods(http.MethodGet)
}

func (s *Server) handleHome() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Welcome to ZROCK API Server")
	}
}

func (s *Server) handleRegister() http.HandlerFunc {
	type registerData struct {
		Email     string `json:"email"`
		Password  string `json:"password"`
		Nickname  string `json:"nickname"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name,omitempty"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &registerData{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			s.error(w, r, http.StatusBadRequest, JSONDecodeError)
			return
		}
		u := &model.User{
			Email:     req.Email,
			Password:  req.Password,
			Nickname:  req.Nickname,
			FirstName: req.FirstName,
			LastName:  req.LastName,
		}
		if err := s.store.User().Create(u); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		u.Sanitize()
		u.Created = nil
		u.LastLogin = nil
		s.respond(w, r, http.StatusCreated, u)
	}
}

func (s *Server) handleLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := &loginData{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			s.error(w, r, http.StatusBadRequest, JSONDecodeError)
			return
		}
		if err := req.validate(); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		u, err := s.store.User().FindByEmailOrNick(req.Username)
		if err != nil || !u.ComparePassword(req.Password) {
			s.error(w, r, http.StatusUnauthorized, ErrUsernameOrPassword)
			return
		}
		u.Sanitize()

		token, err := CreateJWT(u.ID)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		respData := struct {
			*model.User
			Token string `json:"token"`
		}{u, token}
		s.respond(w, r, http.StatusOK, respData)
	}
}

func (s *Server) handlePrivate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Welcome to ZROCK API Server private page")
	}
}

func (s *Server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	if reflect.TypeOf(err).String() == "*errors.errorString" {
		s.respond(w, r, code, map[string]string{"error": err.Error()})
	} else {
		s.respond(w, r, code, map[string]interface{}{"error": err})
	}
}

func (s *Server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			s.logger.Error("Unable to encode server response!")
		}
	}
}

// Middlewares

func (s *Server) authRequiredJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			s.error(w, r, http.StatusUnauthorized, errors.New("not authenticated"))
			return
		}
		authHeaderSplit := strings.Split(authHeader, " ")
		if len(authHeaderSplit) != 2 {
			s.error(w, r, http.StatusUnauthorized, errors.New("not authenticated"))
			return
		}
		userID, err := ParseJWT(authHeaderSplit[1])
		if err != nil {
			s.error(w, r, http.StatusForbidden, errors.New("incorrect token"))
			return
		}
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), contextKeyUserID, userID)))
	})
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
