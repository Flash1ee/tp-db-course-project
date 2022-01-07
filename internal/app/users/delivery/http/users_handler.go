package users_handler

import (
	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"net/http"
	"tp-db-project/internal/app"
	mw "tp-db-project/internal/app/middlewares"
	"tp-db-project/internal/app/users"
	"tp-db-project/internal/app/users/models"
	"tp-db-project/internal/app/users/users_usecase"
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
	utilitiesMiddleware := mw.NewUtilitiesMiddleware(h.logger)
	middlewares := alice.New(context.ClearHandler, utilitiesMiddleware.UpgradeLogger, utilitiesMiddleware.CheckPanic)
	h.router.Get("/api/user/:nickname/profile", middlewares.ThenFunc(h.GetProfileHandler))
	h.router.Post("/api/user/:nickname/profile", middlewares.ThenFunc(h.UpdateProfileHandler))
	h.router.Post("/api/user/:nickname/create", middlewares.ThenFunc(h.CreateUserHandler))

	return h
}
func (h *UsersHandler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	req := &models.User{}
	params, ok := r.Context().Value("params").(httprouter.Params)
	if !ok || len(params) > 1 || params.ByName("nickname") == "" {
		h.Error(w, r, http.StatusBadRequest, InvalidArgument)
		return
	}
	nickname := params.ByName("nickname")

	if err := h.GetRequestBody(w, r, req); err != nil {
		h.Error(w, r, http.StatusBadRequest, InvalidBody)
		return
	}
	req.Nickname = nickname
	if u, err := h.usecase.GetUser(req.Nickname, req.Email); len(u) != 0 || err != nil {
		if err != nil {
			h.UsecaseError(w, r, err, CodeByErrorPost)
			return
		}
		h.Respond(w, r, http.StatusConflict, u)
		return
	}
	user, err := h.usecase.CreateUser(req)
	if err != nil {
		h.UsecaseError(w, r, err, CodeByErrorPost)
		return
	}
	h.Respond(w, r, http.StatusCreated, *user)
}

func (h *UsersHandler) GetProfileHandler(w http.ResponseWriter, r *http.Request) {
	params, ok := r.Context().Value("params").(httprouter.Params)
	if !ok || len(params) > 1 || params.ByName("nickname") == "" {
		h.Error(w, r, http.StatusBadRequest, InvalidArgument)
		return
	}
	nickname := params.ByName("nickname")

	user, err := h.usecase.GetUserByNickname(nickname)
	if err != nil {
		h.UsecaseError(w, r, err, CodeByErrorGet)
		return
	}
	h.Respond(w, r, http.StatusOK, *user)
}
func (h *UsersHandler) UpdateProfileHandler(w http.ResponseWriter, r *http.Request) {
	req := &models.User{}
	params, ok := r.Context().Value("params").(httprouter.Params)
	if !ok || len(params) > 1 || params.ByName("nickname") == "" {
		h.Error(w, r, http.StatusBadRequest, InvalidArgument)
		return
	}
	nickname := params.ByName("nickname")

	if err := h.GetRequestBody(w, r, req); err != nil {
		h.Error(w, r, http.StatusBadRequest, InvalidBody)
		return
	}
	req.Nickname = nickname
	user, err := h.usecase.GetUserByNickname(req.Nickname)
	if err != nil {
		h.UsecaseError(w, r, err, CodeByErrorPost)
		return
	}
	if user == nil {
		h.Error(w, r, http.StatusNotFound, NotFound)
		return
	}
	u, err := h.usecase.UpdateUser(req)
	if err != nil {
		if errors.Cause(err).(*app.GeneralError).Err == users_usecase.ConstraintError {
			h.Error(w, r, http.StatusConflict, ConflictErr)
			return

		}
		h.UsecaseError(w, r, err, CodeByErrorPost)
		return
	}
	u.Nickname = nickname
	h.Respond(w, r, http.StatusOK, *u)
}
