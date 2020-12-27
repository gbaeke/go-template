package api

import (
	"fmt"
	"net/http"
)

func (s *Server) indexHandler(w http.ResponseWriter, r *http.Request) {
	s.logger.Infow("serving index")

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, s.config.Welcome)
}
