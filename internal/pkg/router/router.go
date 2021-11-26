package router

import (
	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type CustomRouter struct {
	*httprouter.Router
}

func NewRouter() *CustomRouter {
	return &CustomRouter{httprouter.New()}
}
func wrapHandler(h http.Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		context.Set(r, "params", ps)
		h.ServeHTTP(w, r)
	}
}
func (r *CustomRouter) Get(path string, handler http.Handler) {
	r.GET(path, wrapHandler(handler))
}
func (r *CustomRouter) Post(path string, handler http.Handler) {
	r.POST(path, wrapHandler(handler))
}
