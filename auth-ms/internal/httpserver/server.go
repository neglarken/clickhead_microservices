package httpserver

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/neglarken/clickhead/auth-ms/config"
	"github.com/neglarken/clickhead/auth-ms/internal/auth"
	"github.com/neglarken/clickhead/auth-ms/internal/handler"
	"github.com/neglarken/clickhead/auth-ms/internal/hasher"
	"github.com/neglarken/clickhead/auth-ms/internal/postgres"
	"github.com/neglarken/clickhead/auth-ms/internal/repo"
	"github.com/neglarken/clickhead/auth-ms/internal/service"
)

type Server struct {
	handler *mux.Router
}

func NewServer(handler mux.Router, sessionStore sessions.Store) *Server {
	s := &Server{
		handler: &handler,
	}
	return s
}

func Start(cfg *config.Config) error {
	db, err := postgres.NewDB(cfg.URL)
	if err != nil {
		return err
	}

	defer db.Close()

	sessionStore := sessions.NewCookieStore([]byte(cfg.SignKey))

	r := repo.NewRepository(db)
	hasher := hasher.NewSHA1Hasher(cfg.Salt)
	manager, err := auth.NewManager(cfg.SignKey)
	if err != nil {
		return err
	}
	serv := service.NewService(*r, hasher, manager, cfg.AccessTokenTTL, cfg.RefreshTokenTTL)
	h := handler.NewHandler(serv, manager)
	s := NewServer(*handler.InitRouter(h), sessionStore)

	h.Logger.Infof("Starting server on %s", cfg.HTTP.Port)
	return http.ListenAndServe(cfg.HTTP.Port, s)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.handler.ServeHTTP(w, r)
}
