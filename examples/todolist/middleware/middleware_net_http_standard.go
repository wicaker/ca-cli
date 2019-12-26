package middleware

import (
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

// StandarMuxMiddleware represent the data-struct for middleware
type StandarMuxMiddleware struct {
	// another stuff , may be needed by middleware
}

// InitStandarMuxMiddleware intialize the middleware
func InitStandarMuxMiddleware() *StandarMuxMiddleware {
	return &StandarMuxMiddleware{}
}

// MiddlewareLogging for logging
func (m *StandarMuxMiddleware) MiddlewareLogging(next http.Handler) http.Handler {
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

// CORS will handle the CORS middleware
func (m *StandarMuxMiddleware) CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		next.ServeHTTP(w, r)
	})
}
