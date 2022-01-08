package service_handler

import (
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"tp-db-project/internal/app/service"
	"tp-db-project/internal/pkg/handler"
	"tp-db-project/internal/pkg/utilits"
)

type ServiceHandler struct {
	router  *mux.Router
	logger  *logrus.Logger
	usecase service.Usecase
	handler.HelpHandlers
	handler.BaseHandler
}

func NewServiceHandler(router *mux.Router, logger *logrus.Logger, uc service.Usecase) *ServiceHandler {
	h := &ServiceHandler{
		router:  router,
		logger:  logger,
		usecase: uc,
		HelpHandlers: handler.HelpHandlers{
			Responder: utilits.Responder{
				LogObject: utilits.NewLogObject(logger),
			},
		},
	}
	h.router.HandleFunc("/api/service/status", h.StatusHandler).Methods("GET")
	h.router.HandleFunc("/api/service/clear", h.ClearHandler).Methods("POST")
	//h.router.Get("/api/service/status", h.StatusHandler)
	//h.router.Post("/api/service/clear", h.ClearHandler)

	return h
}
func (h *ServiceHandler) StatusHandler(w http.ResponseWriter, r *http.Request) {
	status, err := h.usecase.Status()
	if err != nil {
		h.UsecaseError(w, r, err, CodeByErrorGet)
		return
	}
	h.Respond(w, r, http.StatusOK, *status)
}

func (h *ServiceHandler) ClearHandler(w http.ResponseWriter, r *http.Request) {
	err := h.usecase.Clear()
	if err != nil {
		h.UsecaseError(w, r, err, CodeByErrorPost)
		return
	}
	w.WriteHeader(http.StatusOK)
}
