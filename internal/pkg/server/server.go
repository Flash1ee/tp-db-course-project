package server

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"tp-db-project/configs"
	"tp-db-project/internal/pkg/handler"
	http_router "tp-db-project/internal/pkg/router"
	"tp-db-project/internal/pkg/utilits"
)

type Server struct {
	handler     handler.Handler
	config      *configs.Config
	logger      *logrus.Logger
	connections utilits.ExpectedConnections
}

func NewServer(config *configs.Config, connections utilits.ExpectedConnections, logger *logrus.Logger) *Server {
	return &Server{
		config:      config,
		connections: connections,
		logger:      logger,
	}
}
func (s *Server) checkConnection() error {
	if err := s.connections.SqlConnection.Ping(); err != nil {
		return fmt.Errorf("Can't check connection to sql with error %v ", err)
	}
	s.logger.Info("Success check connection to sql db")
	return nil
}

func (s *Server) Start() error {
	if err := s.checkConnection(); err != nil {
		return err
	}

	router := http_router.NewRouter()

	s.logger.Info("Server start")
	return http.ListenAndServe(s.config.BindAddr, router)
}
