package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/neglarken/clickhead/some-ms/internal/repo"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	repo   *repo.Repository
	Logger *logrus.Logger
}

func NewHandler(repo *repo.Repository) *Handler {
	return &Handler{
		repo:   repo,
		Logger: logrus.New(),
	}
}

func InitRouter(h *Handler) *mux.Router {
	router := mux.NewRouter()

	router.Use(h.logRequest)

	router.HandleFunc("/items/", h.CreateItem()).Methods("POST")
	router.HandleFunc("/items/", h.EditItem()).Methods("PUT")
	router.HandleFunc("/items/", h.DeleteItem()).Methods("DELETE")
	router.HandleFunc("/items/", h.GetItems()).Methods("GET")

	return router
}

func (h *Handler) error(w http.ResponseWriter, r *http.Request, httpCode int, err error) {
	h.respond(w, r, httpCode, map[string]string{"error": err.Error()})
}

func (h *Handler) respond(w http.ResponseWriter, r *http.Request, httpCode int, data interface{}) {
	w.WriteHeader(httpCode)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
