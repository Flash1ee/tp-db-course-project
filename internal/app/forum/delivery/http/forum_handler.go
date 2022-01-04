package forum_handler

import (
	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"tp-db-project/internal/app/forum"
	models_forum "tp-db-project/internal/app/forum/models"
	forum_usecase "tp-db-project/internal/app/forum/usecase"
	mw "tp-db-project/internal/app/middlewares"
	"tp-db-project/internal/app/models"
	"tp-db-project/internal/pkg/handler"
	"tp-db-project/internal/pkg/router"
	"tp-db-project/internal/pkg/utilits"
)

const (
	defaultLimit = 100
)

type ForumHandler struct {
	router  *router.CustomRouter
	logger  *logrus.Logger
	usecase forum.Usecase
	handler.HelpHandlers
	handler.BaseHandler
}

func NewForumHandler(router *router.CustomRouter, logger *logrus.Logger, uc forum.Usecase) *ForumHandler {
	h := &ForumHandler{
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
	h.router.Get("/forum/:slug/details", middlewares.ThenFunc(h.GetForumInfo))
	h.router.Get("/forum/:slug/users", middlewares.ThenFunc(h.ForumUsers))
	h.router.Get("/forum/:slug/threads", middlewares.ThenFunc(h.ForumThreads))

	h.router.Post("/forum/:slug/create", middlewares.ThenFunc(h.CreateForumThreads))
	h.router.Post("/forum/create", middlewares.ThenFunc(h.CreateForum))

	return h
}

func (h *ForumHandler) GetForumInfo(w http.ResponseWriter, r *http.Request) {
	params, ok := r.Context().Value("params").(httprouter.Params)
	if !ok || len(params) > 1 || params.ByName("slug") == "" {
		h.Error(w, r, http.StatusBadRequest, InvalidArgument)
		return
	}

	slug := params.ByName("slug")
	res, err := h.usecase.GetForumBySlag(slug)
	if err != nil {
		h.UsecaseError(w, r, err, CodeByErrorGet)
		return
	}
	h.Respond(w, r, http.StatusOK, *res)
}

func (h *ForumHandler) ForumUsers(w http.ResponseWriter, r *http.Request) {
	params, ok := r.Context().Value("params").(httprouter.Params)
	if !ok || len(params) > 1 || params.ByName("slug") == "" {
		h.Error(w, r, http.StatusBadRequest, InvalidArgument)
		return
	}
	slug := params.ByName("slug")
	limit, since, desc := params.ByName("limit"), params.ByName("since"), params.ByName("desc")
	var err error
	var limitInt, sinceInt int
	var descBool bool

	if limit == "" {
		limitInt = defaultLimit
	} else {
		limitInt, err = strconv.Atoi(limit)
		if err != nil {
			h.Error(w, r, http.StatusBadRequest, InvalidLimitArgument)
			return
		}
	}

	pag := &models.Pagination{Limit: int64(limitInt)}

	if since == "" {
		sinceInt = -1
	} else {
		sinceInt, err = strconv.Atoi(since)
		if err != nil {
			h.Error(w, r, http.StatusBadRequest, InvalidLimitArgument)
			return
		}
	}

	if desc != "true" && desc != "false" && desc != "" {
		h.Error(w, r, http.StatusBadRequest, InvalidDescArgument)
		return
	} else {
		if desc == "" || desc == "false" {
			descBool = false
		} else {
			descBool = true
		}
	}

	res, err := h.usecase.GetForumUsers(slug, sinceInt, descBool, pag)
	if err != nil {
		h.UsecaseError(w, r, err, CodeByErrorGet)
		return
	}
	h.Respond(w, r, http.StatusOK, res)
}

func (h *ForumHandler) ForumThreads(w http.ResponseWriter, r *http.Request) {
	params, ok := r.Context().Value("params").(httprouter.Params)
	if !ok || len(params) > 1 || params.ByName("slug") == "" {
		h.Error(w, r, http.StatusBadRequest, InvalidArgument)
		return
	}
	slug := params.ByName("slug")
	limit, since, desc := params.ByName("limit"), params.ByName("since"), params.ByName("desc")
	var err error
	var limitInt, sinceInt int
	var descBool bool

	if limit == "" {
		limitInt = defaultLimit
	} else {
		limitInt, err = strconv.Atoi(limit)
		if err != nil {
			h.Error(w, r, http.StatusBadRequest, InvalidLimitArgument)
			return
		}
	}

	pag := &models.Pagination{Limit: int64(limitInt)}

	if since == "" {
		sinceInt = -1
	} else {
		sinceInt, err = strconv.Atoi(since)
		if err != nil {
			h.Error(w, r, http.StatusBadRequest, InvalidLimitArgument)
			return
		}
	}

	if desc != "true" && desc != "false" && desc != "" {
		h.Error(w, r, http.StatusBadRequest, InvalidDescArgument)
		return
	} else {
		if desc == "" || desc == "false" {
			descBool = false
		} else {
			descBool = true
		}
	}

	res, err := h.usecase.GetForumThreads(slug, sinceInt, descBool, pag)
	if err != nil {
		h.UsecaseError(w, r, err, CodeByErrorGet)
		return
	}
	h.Respond(w, r, http.StatusOK, res)
}

func (h *ForumHandler) CreateForum(w http.ResponseWriter, r *http.Request) {
	req := &models_forum.RequestCreateForum{}
	if err := h.GetRequestBody(w, r, req); err != nil {
		h.Error(w, r, http.StatusBadRequest, InvalidBody)
		return
	}

	if f, err := h.usecase.Create(req); err != nil {
		if err == forum_usecase.AlreadyExists {
			h.Respond(w, r, http.StatusConflict, *f)
			return
		}
		h.UsecaseError(w, r, err, CodeByErrorPost)
		return
	}

	h.Respond(w, r, http.StatusCreated, models_forum.Forum{
		Title:         req.Title,
		UsersNickname: req.User,
		Slug:          req.Slug,
		Posts:         0,
		Threads:       0,
	})
}
func (h *ForumHandler) CreateForumThreads(w http.ResponseWriter, r *http.Request) {
	params, ok := r.Context().Value("params").(httprouter.Params)
	if !ok || len(params) > 1 || params.ByName("slug") == "" {
		h.Error(w, r, http.StatusBadRequest, InvalidArgument)
		return
	}
	slug := params.ByName("slug")

	req := &models_forum.RequestCreateThread{}
	if err := h.GetRequestBody(w, r, req); err != nil {
		h.Error(w, r, http.StatusBadRequest, InvalidBody)
		return
	}

	f, err := h.usecase.CreateThread(slug, req)
	if err != nil {
		if err == forum_usecase.AlreadyExists {
			h.Respond(w, r, http.StatusConflict, *f)
			return
		}
		h.UsecaseError(w, r, err, CodeByErrorPost)
		return
	}

	h.Respond(w, r, http.StatusCreated, *f)
}
