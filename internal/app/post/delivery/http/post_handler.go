package post_handler

import (
	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	mw "tp-db-project/internal/app/middlewares"
	"tp-db-project/internal/app/post"
	models2 "tp-db-project/internal/app/post/models"
	"tp-db-project/internal/pkg/handler"
	"tp-db-project/internal/pkg/router"
	"tp-db-project/internal/pkg/utilits"
)

type PostHandler struct {
	router  *router.CustomRouter
	logger  *logrus.Logger
	usecase post.Usecase
	handler.HelpHandlers
	handler.BaseHandler
}

func NewPostHandler(router *router.CustomRouter, logger *logrus.Logger, uc post.Usecase) *PostHandler {
	h := &PostHandler{
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
	h.router.Get("/post/:id/details", middlewares.ThenFunc(h.GetPost))
	h.router.Post("/post/:id/details", middlewares.ThenFunc(h.UpdatePost))

	return h
}

func (h *PostHandler) GetPost(w http.ResponseWriter, r *http.Request) {
	params, ok := r.Context().Value("params").(httprouter.Params)
	if !ok || len(params) > 1 || params.ByName("id") == "" {
		h.Error(w, r, http.StatusBadRequest, InvalidArgument)
		return
	}

	idInt, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		h.Error(w, r, http.StatusBadRequest, InvalidArgument)
		return
	}
	related := params.ByName("related")
	if related != "" {
		if related != "user" && related != "forum" && related != "thread" {
			h.Error(w, r, http.StatusBadRequest, InvalidParamRelated)
			return
		}
	}

	res, err := h.usecase.GetPost(int64(idInt), related)
	if err != nil {
		h.UsecaseError(w, r, err, CodeByErrorGet)
		return
	}
	h.Respond(w, r, http.StatusOK, *res)
}

func (h *PostHandler) UpdatePost(w http.ResponseWriter, r *http.Request) {
	params, ok := r.Context().Value("params").(httprouter.Params)
	if !ok || len(params) > 1 || params.ByName("id") == "" {
		h.Error(w, r, http.StatusBadRequest, InvalidArgument)
		return
	}
	idInt, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		h.Error(w, r, http.StatusBadRequest, InvalidArgument)
		return
	}
	req := &models2.RequestUpdateMessage{}
	if err := h.GetRequestBody(w, r, req); err != nil {
		h.Error(w, r, http.StatusBadRequest, InvalidBody)
		return
	}
	res, err := h.usecase.UpdatePost(int64(idInt), req)
	if err != nil {
		h.UsecaseError(w, r, err, CodeByErrorPost)
		return
	}
	h.Respond(w, r, http.StatusOK, *res)
}
