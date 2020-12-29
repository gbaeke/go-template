package api

import (
	"time"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

//NewMockServer returns server for testing
func NewMockServer() *Server {
	config := &Config{
		Welcome: "MockHello",
		Port:    9999,
		Log:     false,
		Timeout: 15 * time.Second,
	}

	logger, _ := zap.NewDevelopment()

	return &Server{
		router: mux.NewRouter(),
		logger: logger.Sugar(),
		config: config,
	}
}
