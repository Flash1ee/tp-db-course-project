package handler

import "net/http"

type MiddlewareFunc func(http.Handler) http.Handler

type BaseHandler struct {
	middlewares []MiddlewareFunc
	HelpHandlers
}

func (h *BaseHandler) Use(mwf ...MiddlewareFunc) {
	for _, fn := range mwf {
		h.middlewares = append(h.middlewares, fn)
	}
}
