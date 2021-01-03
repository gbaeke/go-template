package api

import (
	"encoding/json"
	"net/http"
)

// @Summary Liveness probe
// @Description Kubernetes uses this as liveness probe
// @Accept json
// @Produces json
// @Router /healthz [get]
// @Success 200 {string} string "ok"
func (s *Server) healthz(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// @Summary Readiness probe
// @Description Kubernetes uses this as readiness probe
// @Accept json
// @Produces json
// @Router /readyz [get]
// @Success 200 {boolean} boolean true
func (s *Server) readyz(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]bool{"ready": true})
}
