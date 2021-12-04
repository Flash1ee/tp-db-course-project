package users_handler

import (
	"github.com/sirupsen/logrus"
	"tp-db-project/internal/app/users"
	"tp-db-project/internal/pkg/handler"
	"tp-db-project/internal/pkg/router"
	"tp-db-project/internal/pkg/utilits"
)

type UsersHandler struct {
	router  *router.CustomRouter
	logger  *logrus.Logger
	usecase users.Usecase
	handler.HelpHandlers
	handler.BaseHandler
}

func NewUsersHandler(router *router.CustomRouter, logger *logrus.Logger, uc users.Usecase) *UsersHandler {
	h := &UsersHandler{
		router:  router,
		logger:  logger,
		usecase: uc,
		HelpHandlers: handler.HelpHandlers{
			Responder: utilits.Responder{
				LogObject: utilits.NewLogObject(logger),
			},
		},
	}

	//h.router.Get("/balance/{user_id}", h.GetBalanceHandler)
	//h.router.HandleFunc("/balance/{user_id}", h.UpdateBalanceHandler, http.MethodPost)
	//h.router.HandleFunc("/transfer", h.TransferMoneyHandler).Methods(http.MethodPost)
	//utilitiesMiddleware := middlewares.NewUtilitiesMiddleware(h.logger)
	//h.router.Use(utilitiesMiddleware.UpgradeLogger)
	//mux.Router.Use()
	//h.router.Use(utilitiesMiddleware.CheckPanic)

	return h
}
