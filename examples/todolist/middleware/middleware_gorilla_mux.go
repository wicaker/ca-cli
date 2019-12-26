package middleware

import (
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

// GorillaMuxMiddleware represent the data-struct for middleware
type GorillaMuxMiddleware struct {
	// another stuff , may be needed by middleware
}

// InitGorillaMuxMiddleware intialize the middleware
func InitGorillaMuxMiddleware() *GorillaMuxMiddleware {
	return &GorillaMuxMiddleware{}
}

// MiddlewareLogging for logging
func (m *GorillaMuxMiddleware) MiddlewareLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		w.Header().Add("Content-Type", "application/json") // actually this will be moved

		log.WithFields(log.Fields{
			"at":     time.Now().Format("2006-01-02 15:04:05"),
			"method": r.Method,
			"uri":    r.RequestURI,
			"ip":     r.RemoteAddr,
		}).Info("incoming request")

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}
