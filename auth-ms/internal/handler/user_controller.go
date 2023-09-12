package handler

import (
	"encoding/json"
	"net/http"

	"github.com/neglarken/clickhead/auth-ms/internal/model"
)

const (
	sessionName = "auth-ms"
)

func (h *Handler) SignUp() http.HandlerFunc {
	type request struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			h.error(w, r, http.StatusBadRequest, err)
			return
		}

		u := model.User{
			Login:               req.Login,
			UnencryptedPassword: req.Password,
		}

		if err := h.service.User.SignUp(&u); err != nil {
			h.error(w, r, http.StatusInternalServerError, err)
			return
		}

		h.respond(w, r, http.StatusCreated, nil)
	}
}

func (h *Handler) SignIn() http.HandlerFunc {
	type request struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			h.error(w, r, http.StatusBadRequest, err)
			return
		}

		res, err := h.service.User.SignIn(&model.User{
			Login:               req.Login,
			UnencryptedPassword: req.Password,
		})
		if err != nil {
			h.error(w, r, http.StatusInternalServerError, err)
			return
		}

		type response struct {
			AccessToken  string `json:"accessToken"`
			RefreshToken string `json:"refreshToken"`
		}

		h.respond(w, r, http.StatusOK, response{AccessToken: res.AccessToken, RefreshToken: res.RefreshToken})
	}
}

func (h *Handler) handleWhoAmI() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.respond(w, r, http.StatusOK, r.Context().Value(ctxKeyUser))
	}
}
