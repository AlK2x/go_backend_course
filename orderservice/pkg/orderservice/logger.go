package orderservice

import (
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func Logger(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		log.WithFields(log.Fields{
			"method":      r.Method,
			"url":         r.URL,
			"remoteAddr":  r.RemoteAddr,
			"userAgent":   r.UserAgent(),
			"requestTime": time.Since(start),
		}).Info("got a new request")
	})
}
