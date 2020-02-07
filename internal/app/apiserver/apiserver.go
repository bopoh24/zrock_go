package apiserver

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/bopoh24/zrock_go/internal/app/model"
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
	config *Config
}

// NewServer returs server instance
func NewServer(config *Config, store store.IfaceStore) *Server {
	srv := &Server{
		router: mux.NewRouter(),
		logger: logrus.New(),
		store:  store,
	}
	if level, err := logrus.ParseLevel(config.LogLevel); err == nil {
		srv.logger.SetLevel(level)
	}
	srv.configureRouter()
	srv.logger.Infof("Server started at %s ... ", config.BindAdd)
	return srv
}

func (s *Server) configureRouter() {
	//allow CORS
	s.router.Use(handlers.CORS(handlers.AllowedOrigins([]string{"*"})))
	s.router.HandleFunc("/", s.handleHome()).Methods(http.MethodGet)
	s.router.HandleFunc("/register", s.handleRegister()).Methods(http.MethodPost)
	s.router.HandleFunc("/login", s.handleLogin()).Methods(http.MethodPost)
	s.router.HandleFunc("/private", s.handlePrivate()).Methods(http.MethodGet)

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
			s.error(w, r, http.StatusBadRequest, err)
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
			s.error(w, r, http.StatusUnprocessableEntity, err)
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
		// ...
	}
}

func (s *Server) handlePrivate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Welcome to ZROCK API Server private page")
	}
}

func (s *Server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{"error": err.Error()})
}

func (s *Server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			s.logger.Error("Unable to encode server response!")
		}
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
