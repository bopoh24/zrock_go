package apiserver

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"

	_ "github.com/bopoh24/zrock_go/api"
	"github.com/bopoh24/zrock_go/internal/app/model"
	"github.com/bopoh24/zrock_go/internal/app/settings"
	"github.com/bopoh24/zrock_go/internal/app/store"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger"
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

	s.router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"), //The url pointing to API definition
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("#swagger-ui"),
	))

	s.router.HandleFunc("/", s.handleHome()).Methods(http.MethodGet)
	authRoutes := s.router.PathPrefix("/api/v1/auth").Subrouter()
	authRoutes.HandleFunc("/login", s.handleLogin()).Methods(http.MethodPost)
	authRoutes.HandleFunc("/register", s.handleRegister()).Methods(http.MethodPost)

	// auth required urls
	api := s.router.PathPrefix("/api/v1").Subrouter()
	api.Use(s.authRequiredJWT)
	api.HandleFunc("/private", s.handlePrivate()).Methods(http.MethodGet)
}

// @Summary Main page
// @Description main page for anonimous access
// @ID main-page
// @Produce  html
// @Success 200
// @Router / [get]
func (s *Server) handleHome() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Welcome to ZROCK API Server")
	}
}

// @Summary Registration
// @Description New user registration method
// @ID registration
// @Param JSON body registerData true "Register data"
// @Accept  json
// @Produce  json
// @Success 201
// @Failure 400 {object} errorResp
// @Router /auth/register [post]
func (s *Server) handleRegister() http.HandlerFunc {
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
		u.Password = ""
		u.Created = nil
		u.LastLogin = nil
		s.respond(w, r, http.StatusCreated, u)
	}
}

// @Summary Login
// @Description User login method
// @ID login
// @Param JSON body loginData true "Login data"
// @Accept  json
// @Produce json
// @Success 200 {object} loginResponseData
// @Failure 400 {object} errorResp
// @Failure 401 {object} errorResp
// @Router /auth/login [post]
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
		if !u.EmailVerified {
			s.error(w, r, http.StatusUnauthorized, ErrEmailNotVerified)
			return
		}
		u.Sanitize()

		token, err := CreateJWT(u.ID)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		s.respond(w, r, http.StatusOK, &loginResponseData{u, token})
	}
}

// @Summary Private page for authenticated users only
// @Description returns user ID in greating
// @Security Bearer
// @ID private-page
// @Accept  json
// @Produce  json
// @Success 200 {string} string
// @Failure 401
// @Router /private [get]
func (s *Server) handlePrivate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(contextKeyUserID)
		io.WriteString(w, fmt.Sprintf("USER #%d: Welcome to ZROCK API Server private page", userID))
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
