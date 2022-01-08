package post_handler

import (
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"strings"
	"tp-db-project/internal/app/post"
	models2 "tp-db-project/internal/app/post/models"
	post_repository "tp-db-project/internal/app/post/repository"
	"tp-db-project/internal/pkg/handler"
	"tp-db-project/internal/pkg/utilits"
)

type PostHandler struct {
	router  *mux.Router
	logger  *logrus.Logger
	usecase post.Usecase
	handler.HelpHandlers
	handler.BaseHandler
}

func NewPostHandler(router *mux.Router, logger *logrus.Logger, uc post.Usecase) *PostHandler {
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
	h.router.HandleFunc("/api/post/{id}/details", h.GetPost).Methods("GET")
	h.router.HandleFunc("/api/post/{id}/details", h.UpdatePost).Methods("POST")
	//h.router.Get("/api/post/:id/details", h.GetPost)
	//h.router.Post("/api/post/:id/details", h.UpdatePost)

	return h
}

func (h *PostHandler) GetPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok || id == "" {
		h.Error(w, r, http.StatusBadRequest, InvalidArgument)
		return
	}
	//params, ok := r.Context().Value("params").(httprouter.Params)
	//if !ok || len(params) > 1 || params.ByName("id") == "" {
	//	h.Error(w, r, http.StatusBadRequest, InvalidArgument)
	//	return
	//}

	//idInt, err := strconv.Atoi(params.ByName("id"))
	idInt, err := strconv.Atoi(id)
	if err != nil {
		h.Error(w, r, http.StatusBadRequest, InvalidArgument)
		return
	}
	related := r.URL.Query().Get("related")
	if related != "" {
		args := strings.Split(related, ",")
		for _, arg := range args {
			if arg != "user" && arg != "forum" && arg != "thread" {
				h.Error(w, r, http.StatusBadRequest, InvalidParamRelated)
				return
			}
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
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok || id == "" {
		h.Error(w, r, http.StatusBadRequest, InvalidArgument)
		return
	}
	//params, ok := r.Context().Value("params").(httprouter.Params)
	//if !ok || len(params) > 1 || params.ByName("id") == "" {
	//	h.Error(w, r, http.StatusBadRequest, InvalidArgument)
	//	return
	//}
	//idInt, err := strconv.Atoi(params.ByName("id"))
	idInt, err := strconv.Atoi(id)
	if err != nil {
		h.Error(w, r, http.StatusBadRequest, InvalidArgument)
		return
	}
	req := &models2.RequestUpdateMessage{}
	if err := h.GetRequestBody(w, r, req); err != nil {
		if err.Error() == "EOF" {
			h.Error(w, r, http.StatusNotFound, post_repository.NotFound)
			return
		}
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
