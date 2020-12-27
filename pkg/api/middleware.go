package api

import (
	"net/http"

	"go.uber.org/zap"
)

func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// info logging
		s.logger.Infow("ran handler",
			zap.String("requestURI", r.RequestURI),
		)
		// Call the next handler
		next.ServeHTTP(w, r)
	})
}
