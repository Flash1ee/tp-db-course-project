package server

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"tp-db-project/configs"
	users_handler "tp-db-project/internal/app/users/delivery/http"
	"tp-db-project/internal/app/users/users_repository"
	"tp-db-project/internal/app/users/users_usecase"
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
	if err := s.connections.SqlConnection.Ping(context.Background()); err != nil {
		return fmt.Errorf("Can't check connection to sql with error %v ", err)
	}
	s.logger.Info("Success check connection to sql db")
	return nil
}

func (s *Server) Start() error {
	if err := s.checkConnection(); err != nil {
		return err
	}

	router := http_router.NewRouter(s.logger)

	userUsecase := users_usecase.NewUsersUsecase(users_repository.NewUsersRepository(s.connections.SqlConnection))
	_ = users_handler.NewUsersHandler(router, s.logger, userUsecase)

	s.logger.Info("Server start")
	return http.ListenAndServe(s.config.BindAddr, router)
}
