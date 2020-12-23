package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

//Config API configuration via viper
type Config struct {
	Welcome string
	Port    int
}

//Server struct
type Server struct {
	config *Config
	logger *zap.SugaredLogger
	router *mux.Router
}

//NewServer creates new server
func NewServer(config *Config, logger *zap.SugaredLogger) (*Server, error) {
	srv := &Server{
		config: config,
		logger: logger,
		router: mux.NewRouter(),
	}

	return srv, nil
}

func (s *Server) indexRoute(w http.ResponseWriter, r *http.Request) {
	s.logger.Infow("serving index")

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, s.config.Welcome)
}

//SetupRoutes sets up routes
func (s *Server) SetupRoutes() {
	s.router.HandleFunc("/", s.indexRoute)
}

//StartServer starts http server
func (s *Server) StartServer() {

	s.SetupRoutes()

	srv := &http.Server{
		Addr:    ":" + fmt.Sprint(s.config.Port),
		Handler: s.router,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	s.logger.Infow("starting web server",
		zap.Int("port", s.config.Port),
	)

	s.logger.Fatal(srv.ListenAndServe())
}
