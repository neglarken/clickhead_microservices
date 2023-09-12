package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/neglarken/clickhead/auth-ms/internal/auth"
	"github.com/neglarken/clickhead/auth-ms/internal/service"
	"github.com/sirupsen/logrus"
)

const (
	ctxKeyUser = "userId"
)

type Handler struct {
	service      *service.Service
	Logger       *logrus.Logger
	tokenManager auth.TokenManager
}

func NewHandler(service *service.Service, tokenManager auth.TokenManager) *Handler {
	return &Handler{
		service:      service,
		Logger:       logrus.New(),
		tokenManager: tokenManager,
	}
}

func InitRouter(h *Handler) *mux.Router {
	router := mux.NewRouter()

	router.Use(h.logRequest)

	router.HandleFunc("/users/sign-up", h.SignUp()).Methods("POST")
	router.HandleFunc("/users/sign-in", h.SignIn()).Methods("POST")

	private := router.PathPrefix("/private").Subrouter()
	private.Use(h.authUser)
	private.HandleFunc("/whoAmI", h.handleWhoAmI()).Methods("GET")

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
