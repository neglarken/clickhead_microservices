package httpserver

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/neglarken/clickhead/some-ms/config"
	"github.com/neglarken/clickhead/some-ms/internal/handler"
	"github.com/neglarken/clickhead/some-ms/internal/postgres"
	"github.com/neglarken/clickhead/some-ms/internal/repo"
)

type Server struct {
	handler *mux.Router
}

func NewServer(handler mux.Router) *Server {
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

	r := repo.NewRepository(db)
	h := handler.NewHandler(r)
	s := NewServer(*handler.InitRouter(h))

	h.Logger.Infof("Starting server on %s", cfg.HTTP.Port)
	return http.ListenAndServe(cfg.HTTP.Port, s)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.handler.ServeHTTP(w, r)
}
