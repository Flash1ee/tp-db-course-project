package router

import (
	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type CustomRouter struct {
	router *httprouter.Router
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
	r.router.GET(path, wrapHandler(handler))
}
func (r *CustomRouter) Post(path string, handler http.Handler) {
	r.router.POST(path, wrapHandler(handler))
}
func (r *CustomRouter) HandleFunc(url string, f func(http.ResponseWriter, *http.Request), method string) {
	r.router.HandlerFunc(method, url, f)
}
