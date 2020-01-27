package zrockapi

import (
	"io"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// Server api
type Server struct {
	router *mux.Router
	logger *logrus.Logger
	config *Config
}

// NewServer returs server instance
func NewServer(config *Config) *Server {
	srv := &Server{
		router: mux.NewRouter(),
		logger: logrus.New(),
	}
	srv.configureRouter()
	return srv
}

func (s *Server) configureRouter() {
	//allow CORS
	s.router.Use(handlers.CORS(handlers.AllowedOrigins([]string{"*"})))
	s.router.HandleFunc("/", s.handleHome()).Methods(http.MethodGet)
}

func (s *Server) handleHome() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Welcome to ZROCK API Server")
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
