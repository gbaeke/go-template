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

	_ "github.com/gbaeke/go-template/pkg/api/docs"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/swaggo/swag"
)

// @title go-template API
// @version 0.1
// @description Go template

// @contact.name Source Code
// @contact.url https://github.com/gbaeke/go-template

// @host localhost:8080
// @BasePath /
// @schemes http https

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
	s.router.Handle("/metrics", promhttp.Handler())
	s.router.HandleFunc("/healthz", s.healthz)
	s.router.HandleFunc("/readyz", s.readyz)
	s.router.HandleFunc("/", s.indexHandler)
	s.router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))
	s.router.HandleFunc("/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		doc, err := swag.ReadDoc()
		if err != nil {
			s.logger.Error("swagger error", zap.Error(err), zap.String("path", "/swagger.json"))
		}
		w.Write([]byte(doc))
	})
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
