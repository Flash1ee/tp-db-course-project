package forum_handler

import (
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"tp-db-project/internal/app/forum"
	models2 "tp-db-project/internal/app/forum/delivery/models"
	models_forum "tp-db-project/internal/app/forum/models"
	forum_usecase "tp-db-project/internal/app/forum/usecase"
	"tp-db-project/internal/app/models"
	"tp-db-project/internal/pkg/handler"
	"tp-db-project/internal/pkg/utilits"
)

const (
	defaultLimit = 100
)

type ForumHandler struct {
	router  *mux.Router
	logger  *logrus.Logger
	usecase forum.Usecase
	handler.HelpHandlers
	handler.BaseHandler
}

func NewForumHandler(router *mux.Router, logger *logrus.Logger, uc forum.Usecase) *ForumHandler {
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
	//h.router.Get("/forum/:slug/details", middlewares.ThenFunc(h.GetForumInfo))
	//h.router.Get("/forum/:slug/users", middlewares.ThenFunc(h.ForumUsers))
	//h.router.Get("/forum/:slug/threads", middlewares.ThenFunc(h.ForumThreads))

	//h.router.Post("/forum/create", middlewares.ThenFunc(h.CreateForum))
	//h.router.Post("/forum/:slug/create", middlewares.ThenFunc(h.CreateForumThreads))

	h.router.HandleFunc("/api/forum/{slug}/details", h.GetForumInfo).Methods("GET")
	h.router.HandleFunc("/api/forum/{slug}/users", h.ForumUsers).Methods("GET")
	h.router.HandleFunc("/api/forum/{slug}/threads", h.ForumThreads).Methods("GET")
	h.router.HandleFunc("/api/forum/create", h.CreateForum).Methods("POST")
	h.router.HandleFunc("/api/forum/{slug}/create", h.CreateForumThreads).Methods("POST")

	return h
}

func (h *ForumHandler) GetForumInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug, ok := vars["slug"]
	if !ok || slug == "" {
		h.Error(w, r, http.StatusBadRequest, InvalidArgument)
		return
	}

	//slug := params.ByName("slug")
	res, err := h.usecase.GetForum(slug)
	if err != nil {
		//if errors.Cause(err).(*app.GeneralError).Err == forum_usecase.ForumNotFound {
		//	h.Respond(w, r, http.StatusOK, nil)
		//	return
		//}
		h.UsecaseError(w, r, err, CodeByErrorGet)
		return
	}
	h.Respond(w, r, http.StatusOK, *res)
}

func (h *ForumHandler) ForumUsers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug, ok := vars["slug"]
	if !ok || slug == "" {
		h.Error(w, r, http.StatusBadRequest, InvalidArgument)
		return
	}
	//params, ok := r.Context().Value("params").(httprouter.Params)
	//if !ok || len(params) > 1 || params.ByName("slug") == "" {
	//	h.Error(w, r, http.StatusBadRequest, InvalidArgument)
	//	return
	//}
	//slug := params.ByName("slug")
	limit := r.URL.Query().Get("limit")
	since := r.URL.Query().Get("since")
	desc := r.URL.Query().Get("desc")
	//limit, since, desc := params.ByName("limit"), params.ByName("since"), params.ByName("desc")
	var err error
	var limitInt int
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

	//if since == "" {
	//	sinceInt = -1
	//} else {
	//	sinceInt, err = strconv.Atoi(since)
	//	if err != nil {
	//		h.Error(w, r, http.StatusBadRequest, InvalidLimitArgument)
	//		return
	//	}
	//}

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

	res, err := h.usecase.GetForumUsers(slug, since, descBool, pag)
	if err != nil {
		h.UsecaseError(w, r, err, CodeByErrorGet)
		return
	}
	h.Respond(w, r, http.StatusOK, res)
}

func (h *ForumHandler) ForumThreads(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug, ok := vars["slug"]
	if !ok || slug == "" {
		h.Error(w, r, http.StatusBadRequest, InvalidArgument)
		return
	}
	//params, ok := r.Context().Value("params").(httprouter.Params)
	//if !ok || len(params) > 1 || params.ByName("slug") == "" {
	//	h.Error(w, r, http.StatusBadRequest, InvalidArgument)
	//	return
	//}
	//slug := params.ByName("slug")
	limit := r.URL.Query().Get("limit")
	since := r.URL.Query().Get("since")
	desc := r.URL.Query().Get("desc")
	//limit, since, desc := params.ByName("limit"), params.ByName("since"), params.ByName("desc")
	//params, ok := r.Context().Value("params").(httprouter.Params)
	//if !ok || len(params) > 1 || params.ByName("slug") == "" {
	//	h.Error(w, r, http.StatusBadRequest, InvalidArgument)
	//	return
	//}
	//slug := params.ByName("slug")
	//limit, since, desc := params.ByName("limit"), params.ByName("since"), params.ByName("desc")
	var err error
	var limitInt int
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

	//if since == "" {
	//	sinceInt = -1
	//} else {
	//	sinceInt, err = strconv.Atoi(since)
	//	if err != nil {
	//		h.Error(w, r, http.StatusBadRequest, InvalidLimitArgument)
	//		return
	//	}
	//}

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

	res, err := h.usecase.GetForumThreads(slug, since, descBool, pag)
	if err != nil {
		//if errors.Cause(err).(*app.GeneralError).Err == forum_usecase.ForumNotFound {
		//	h.Respond(w, r, http.StatusOK, nil)
		//	return
		//}
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

	h.Respond(w, r, http.StatusCreated, models2.ResponseForum{
		Title:         req.Title,
		UsersNickname: req.User,
		Slug:          req.Slug,
	})
}
func (h *ForumHandler) CreateForumThreads(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slugForum, ok := vars["slug"]
	if !ok || slugForum == "" {
		h.Error(w, r, http.StatusBadRequest, InvalidArgument)
		return
	}
	//params, ok := r.Context().Value("params").(httprouter.Params)
	//if !ok || len(params) > 1 || params.ByName("slug") == "" {
	//	h.Error(w, r, http.StatusBadRequest, InvalidArgument)
	//	return
	//}
	//slug := params.ByName("slug")

	req := &models_forum.RequestCreateThread{}
	if err := h.GetRequestBody(w, r, req); err != nil {
		h.Error(w, r, http.StatusBadRequest, InvalidBody)
		return
	}
	//flag := false
	if req.Forum == "" {
		req.Forum = slugForum
		//flag = true
	}
	//fmt.Printf("forum slug = %s; thread slug = %s\n", req.Forum, req.Slug)
	f, err := h.usecase.CreateThread(req)
	if err != nil {
		if err == forum_usecase.AlreadyExists {
			h.Respond(w, r, http.StatusConflict, *f)
			return
		}
		h.UsecaseError(w, r, err, CodeByErrorPost)
		return
	}

	//if flag {
	//	f.Forum = InvertCase(f.Forum)
	//}

	h.Respond(w, r, http.StatusCreated, *f)
}
