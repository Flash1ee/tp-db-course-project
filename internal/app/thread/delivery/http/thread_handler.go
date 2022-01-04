package thread_handler

import (
	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"github.com/sirupsen/logrus"
	"net/http"
	mw "tp-db-project/internal/app/middlewares"
	"tp-db-project/internal/app/thread"
	"tp-db-project/internal/app/thread/models"
	models2 "tp-db-project/internal/app/vote/models"
	"tp-db-project/internal/pkg/handler"
	"tp-db-project/internal/pkg/router"
	"tp-db-project/internal/pkg/utilits"
)

type ThreadHandler struct {
	router  *router.CustomRouter
	logger  *logrus.Logger
	usecase thread.Usecase
	handler.HelpHandlers
	handler.BaseHandler
}

func NewThreadHandler(router *router.CustomRouter, logger *logrus.Logger, uc thread.Usecase) *ThreadHandler {
	h := &ThreadHandler{
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
	h.router.Get("/thread/:slug_or_id/details", middlewares.ThenFunc(h.ThreadInfo))
	h.router.Get("/thread/:slug_or_id/posts", middlewares.ThenFunc(h.ThreadPosts))

	h.router.Post("/thread/:slug_or_id/details", middlewares.ThenFunc(h.UpdateThread))
	h.router.Post("/thread/:slug_or_id/vote", middlewares.ThenFunc(h.VoteThread))

	return h
}

func (h *ThreadHandler) ThreadInfo(w http.ResponseWriter, r *http.Request) {
	params, ok := r.Context().Value("params").(httprouter.Params)
	if !ok || len(params) > 1 || params.ByName("slug_or_id") == "" {
		h.Error(w, r, http.StatusBadRequest, InvalidArgument)
		return
	}
	slugOrID := params.ByName("slug_or_id")
	res, err := h.usecase.GetThreadInfo(slugOrID)
	if err != nil {
		h.UsecaseError(w, r, err, CodeByErrorGet)
		return
	}
	h.Respond(w, r, http.StatusOK, *res)

}
func (h *ThreadHandler) ThreadPosts(w http.ResponseWriter, r *http.Request) {
	//params, ok := r.Context().Value("params").(httprouter.Params)
	//if !ok || len(params) > 1 || params.ByName("slug_or_id") == "" {
	//	h.Error(w, r, http.StatusBadRequest, InvalidArgument)
	//	return
	//}
	//slugOrID := params.ByName("slug_or_id")
	h.Respond(w, r, http.StatusOK, nil)

}

func (h *ThreadHandler) UpdateThread(w http.ResponseWriter, r *http.Request) {
	params, ok := r.Context().Value("params").(httprouter.Params)
	if !ok || len(params) > 1 || params.ByName("slug_or_id") == "" {
		h.Error(w, r, http.StatusBadRequest, InvalidArgument)
		return
	}
	slugOrID := params.ByName("slug_or_id")
	req := &models.RequestUpdateThread{}
	if err := h.GetRequestBody(w, r, req); err != nil {
		h.Error(w, r, http.StatusBadRequest, InvalidBody)
		return
	}
	res, err := h.usecase.UpdateThread(slugOrID, req)
	if err != nil {
		h.UsecaseError(w, r, err, CodeByErrorPost)
		return
	}
	h.Respond(w, r, http.StatusOK, *res)
}
func (h *ThreadHandler) VoteThread(w http.ResponseWriter, r *http.Request) {
	params, ok := r.Context().Value("params").(httprouter.Params)
	if !ok || len(params) > 1 || params.ByName("slug_or_id") == "" {
		h.Error(w, r, http.StatusBadRequest, InvalidArgument)
		return
	}
	slugOrID := params.ByName("slug_or_id")
	req := &models2.RequestVoteUpdate{}
	if err := h.GetRequestBody(w, r, req); err != nil {
		h.Error(w, r, http.StatusBadRequest, InvalidBody)
		return
	}
	if req.Voice != -1 && req.Voice != 1 {
		h.Error(w, r, http.StatusBadRequest, InvalidVoice)
		return
	}
	_, err := h.usecase.UpdateVoice(slugOrID, req)
	if err != nil {
		h.UsecaseError(w, r, err, CodeByErrorPost)
		return
	}
	threadInfo, err := h.usecase.GetThreadInfo(slugOrID)
	if err != nil {
		h.UsecaseError(w, r, err, CodeByErrorGet)
		return
	}
	h.Respond(w, r, http.StatusOK, *threadInfo)

}
