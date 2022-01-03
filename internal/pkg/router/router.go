package router

import (
	go_context "context"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
	"net/http"
)

type CustomRouter struct {
	router *httprouter.Router
	logger *logrus.Logger
}

func NewRouter(logger *logrus.Logger) *CustomRouter {
	return &CustomRouter{router: httprouter.New(), logger: logger}
}
func wrapHandler(h http.Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		r = r.WithContext(go_context.WithValue(r.Context(), "params", ps))
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
func (rout *CustomRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rout.logger.Info(r)
	rout.router.ServeHTTP(w, r)
}
