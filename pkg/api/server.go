package api

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

//Config API configuration via viper
type Config struct {
	Welcome string
	Port    int
	Log     bool
	Timeout time.Duration
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

//SetupRoutes sets up routes
func (s *Server) setupRoutes() {
	s.router.HandleFunc("/healthz", s.healthz)
	s.router.HandleFunc("/readyz", s.readyz)
	s.router.HandleFunc("/", s.indexHandler)
}

func (s *Server) setupMiddlewares() {
	if s.config.Log {
		// only log requests when --log is set
		s.router.Use(s.loggingMiddleware)
	}
}

//StartServer starts http server
func (s *Server) StartServer() {

	s.setupRoutes()
	s.setupMiddlewares()

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

	//graceful shutdown - run server in goroutine and handle SIGINT & SIGTERM
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			s.logger.Infow("server stopped",
				zap.Error(err),
			)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	// block wait for signal
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), s.config.Timeout)
	defer cancel()
	srv.Shutdown(ctx)

	s.logger.Infow("server shutting down")
	os.Exit(0)
}
