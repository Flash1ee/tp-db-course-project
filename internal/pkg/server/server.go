package server

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"tp-db-project/configs"
	forum_handler "tp-db-project/internal/app/forum/delivery/http"
	"tp-db-project/internal/app/forum/repository/postgresql"
	forum_usecase "tp-db-project/internal/app/forum/usecase"
	post_handler "tp-db-project/internal/app/post/delivery/http"
	post_postgresql "tp-db-project/internal/app/post/repository/postgresql"
	post_usecase "tp-db-project/internal/app/post/usecase"
	service_handler "tp-db-project/internal/app/service/delivery/http"
	"tp-db-project/internal/app/service/repository/postgresql"
	"tp-db-project/internal/app/service/service_usecase"
	thread_handler "tp-db-project/internal/app/thread/delivery/http"
	thread_postgresql "tp-db-project/internal/app/thread/repository/postgresql"
	thread_usecase "tp-db-project/internal/app/thread/usecase"
	users_handler "tp-db-project/internal/app/users/delivery/http"
	"tp-db-project/internal/app/users/users_repository/postgresql"
	"tp-db-project/internal/app/users/users_usecase"
	vote_postgresql "tp-db-project/internal/app/vote/repository/postgresql"
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
	routerForum := http_router.NewMuxRouter(s.logger)

	usersRepo := users_postgresql.NewUsersRepository(s.connections.SqlConnection)
	forumRepo := forum_postgresql.NewForumRepository(s.connections.SqlConnection)
	threadRepo := thread_postgresql.NewThreadRepository(s.connections.SqlConnection)
	postRepo := post_postgresql.NewPostRepository(s.connections.SqlConnection)
	serviceRepo := service_postgresql.NewServiceRepository(s.connections.SqlConnection)
	voteRepo := vote_postgresql.NewVoteRepository(s.connections.SqlConnection)

	userUsecase := users_usecase.NewUsersUsecase(usersRepo)
	postUsecase := post_usecase.NewPostUsecase(postRepo)
	forumUsecase := forum_usecase.NewForumUsecase(forumRepo, usersRepo, threadRepo)
	serviceUsecase := service_usecase.NewServiceUsecase(serviceRepo)
	threadUsecase := thread_usecase.NewThreadUsecase(threadRepo, voteRepo, postRepo, usersRepo)

	_ = users_handler.NewUsersHandler(router, s.logger, userUsecase)
	_ = post_handler.NewPostHandler(router, s.logger, postUsecase)
	_ = forum_handler.NewForumHandler(routerForum, s.logger, forumUsecase)
	_ = service_handler.NewServiceHandler(router, s.logger, serviceUsecase)
	_ = thread_handler.NewThreadHandler(router, s.logger, threadUsecase)

	s.logger.Info("Server start")
	server := http.NewServeMux()
	server.Handle("/forum/", routerForum)
	server.Handle("/thread/", router)
	server.Handle("/service/", router)
	server.Handle("/user/", router)
	server.Handle("/post/", router)

	return http.ListenAndServe(s.config.BindAddr, server)
}
