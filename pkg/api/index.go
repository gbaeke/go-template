package api

import (
	"fmt"
	"net/http"
)

func (s *Server) indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, s.config.Welcome)
}
