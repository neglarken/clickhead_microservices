package handler

import (
	"net/http"
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
