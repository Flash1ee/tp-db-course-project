package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	_ "net/http/pprof"
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


func (s *Server) Start() error {
	//router := http_router.NewRouter(s.logger)
	//routerForum := http_router.NewMuxRouter(s.logger)
	//
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
	muxRouter := mux.NewRouter()

	_ = users_handler.NewUsersHandler(muxRouter, s.logger, userUsecase)
	_ = post_handler.NewPostHandler(muxRouter, s.logger, postUsecase)
	_ = forum_handler.NewForumHandler(muxRouter, s.logger, forumUsecase)
	_ = service_handler.NewServiceHandler(muxRouter, s.logger, serviceUsecase)
	_ = thread_handler.NewThreadHandler(muxRouter, s.logger, threadUsecase)

	s.logger.Info("Server start")
	//server := http.NewServeMux()
	//server.Handle("/api/forum/", routerForum)
	//server.Handle("/api/thread/", router)
	//server.Handle("/api/service/", router)
	//server.Handle("/api/user/", router)
	//server.Handle("/api/post/", router)

	//muxRouter.HandleFunc("/debug/pprof/", pprof.Index).Methods("GET")
	//muxRouter.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline).Methods("GET")
	//muxRouter.HandleFunc("/debug/pprof/profile", pprof.Profile).Methods("GET")
	//muxRouter.HandleFunc("/debug/pprof/symbol", pprof.Symbol).Methods("GET")
	//muxRouter.HandleFunc("/debug/pprof/trace", pprof.Trace).Methods("GET")
	//server := http.Server{
	//	Addr:        s.config.BindAddr,
	//	Handler:      muxRouter,
	//	ReadTimeout:  60 * time.Second,
	//	WriteTimeout: 60 * time.Second,
	//}
	//return server.ListenAndServe()
	checkAliveRouter := mux.NewRouter()

	checkAliveRouter.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Server One")
    })

	go http.ListenAndServe(":81", checkAliveRouter)


	return http.ListenAndServe(s.config.BindAddr, muxRouter)
}
