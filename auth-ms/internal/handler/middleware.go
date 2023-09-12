package handler

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

func (h *Handler) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := h.Logger.WithFields(logrus.Fields{
			"remote_addr": r.RemoteAddr,
		})
		logger.Infof("started %s %s", r.Method, r.RequestURI)
		start := time.Now()
		rw := &responseWriter{w, http.StatusOK}
		next.ServeHTTP(rw, r)
		logger.Infof(
			"completed with %d %s in %v",
			rw.code,
			http.StatusText(rw.code),
			time.Now().Sub(start),
		)
	})
}

func (h *Handler) authUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := h.parseAuthHeader(r)
		if err != nil {
			h.error(w, r, http.StatusUnauthorized, err)
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyUser, id)))
	})
}

func (h *Handler) parseAuthHeader(r *http.Request) (string, error) {
	header := r.Header.Get("Authorization")
	if header == "" {
		return "", errors.New("empty auth header")
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return "", errors.New("invalid auth header")
	}

	if len(headerParts[1]) == 0 {
		return "", errors.New("token is empty")
	}

	return h.tokenManager.Parse(headerParts[1])
}
