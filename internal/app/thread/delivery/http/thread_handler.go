package thread_handler

import (
	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"tp-db-project/internal/app"
	mw "tp-db-project/internal/app/middlewares"
	models3 "tp-db-project/internal/app/models"
	"tp-db-project/internal/app/thread"
	"tp-db-project/internal/app/thread/models"
	thread_usecase "tp-db-project/internal/app/thread/usecase"
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
	h.router.Get("/api/thread/:slug_or_id/details", middlewares.ThenFunc(h.ThreadInfo))
	h.router.Get("/api/thread/:slug_or_id/posts", middlewares.ThenFunc(h.ThreadPosts))

	h.router.Post("/api/thread/:slug_or_id/details", middlewares.ThenFunc(h.UpdateThread))
	h.router.Post("/api/thread/:slug_or_id/vote", middlewares.ThenFunc(h.VoteThread))
	h.router.Post("/api/thread/:slug_or_id/create", middlewares.ThenFunc(h.CreatePosts))
	return h
}
func (h *ThreadHandler) CreatePosts(w http.ResponseWriter, r *http.Request) {
	params, ok := r.Context().Value("params").(httprouter.Params)
	if !ok || len(params) > 1 || params.ByName("slug_or_id") == "" {
		h.Error(w, r, http.StatusBadRequest, InvalidArgument)
		return
	}
	slugOrID := params.ByName("slug_or_id")
	var posts []*models.RequestNewPost

	if err := h.GetRequestBody(w, r, &posts); err != nil {
		h.Error(w, r, http.StatusBadRequest, InvalidBody)
		return
	}
	//if len(posts) == 0 {
	//	h.Respond(w, r, http.StatusCreated, []models4.ResponsePost{})
	//	return
	//}
	res, err := h.usecase.CreatePosts(slugOrID, posts)
	if err != nil {
		if errors.Cause(err).(*app.GeneralError).Err == thread_usecase.CreatedEmpty {
			h.Respond(w, r, http.StatusCreated, []interface{}{})
			return
		}
		h.UsecaseError(w, r, err, CodeByErrorPost)
		return
	}

	h.Respond(w, r, http.StatusCreated, res)
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
	params, ok := r.Context().Value("params").(httprouter.Params)
	if !ok || len(params) > 1 || params.ByName("slug_or_id") == "" {
		h.Error(w, r, http.StatusBadRequest, InvalidArgument)
		return
	}
	var err error
	slugOrID := params.ByName("slug_or_id")
	limit := r.URL.Query().Get("limit")
	since := r.URL.Query().Get("since")
	sort := r.URL.Query().Get("sort")
	desc := r.URL.Query().Get("desc")

	limitInt := 100
	sinceInt := 0
	descBool := false

	if limit != "" {
		limitInt, err = strconv.Atoi(limit)
		if err != nil {
			h.Error(w, r, http.StatusBadRequest, InvalidArgument)
			return
		}
	}
	if since == "" {
		sinceInt = -1
	} else {
		sinceInt, err = strconv.Atoi(since)
		if err != nil {
			h.Error(w, r, http.StatusBadRequest, InvalidArgument)
			return
		}
	}
	//if since != "" {
	//	if sinceInt, err = strconv.Atoi(since); err != nil {
	//		h.Error(w, r, http.StatusBadRequest, InvalidArgument)
	//		return
	//	}
	//}
	if sort != "flat" && sort != "tree" && sort != "parent_tree" &&
		sort != "" {
		h.Error(w, r, http.StatusBadRequest, InvalidArgument)
		return
	}

	if desc == "true" {
		descBool = true
	} else if desc != "false" && desc != "" {
		h.Error(w, r, http.StatusBadRequest, InvalidArgument)
		return
	}

	res, err := h.usecase.GetPostsBySort(slugOrID, sort, int64(sinceInt), descBool, &models3.Pagination{
		Limit: int64(limitInt),
	})
	if err != nil {
		h.UsecaseError(w, r, err, CodeByErrorGet)
		return
	}
	h.Respond(w, r, http.StatusOK, res)
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
	th, err := h.usecase.UpdateVoice(slugOrID, req)
	if err != nil {
		h.UsecaseError(w, r, err, CodeByErrorPost)
		return
	}
	//threadInfo, err := h.usecase.GetThreadInfo(slugOrID)
	//if err != nil {
	//	h.UsecaseError(w, r, err, CodeByErrorGet)
	//	return
	//}
	h.Respond(w, r, http.StatusOK, *th)

}
