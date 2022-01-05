package router

import (
	go_context "context"
	"github.com/gorilla/mux"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
	"net/http"
)

type Router interface {
	HandleFunc(url string, h http.Handler, method string)
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type MuxCustomRouter struct {
	router *mux.Router
	logger *logrus.Logger
}

func NewMuxRouter(logger *logrus.Logger) *MuxCustomRouter {
	return &MuxCustomRouter{router: mux.NewRouter(), logger: logger}
}
func (r *MuxCustomRouter) HandleFunc(url string, h http.Handler, method string) {
	r.router.Handle(url, h).Methods(method)
}
func (rout *MuxCustomRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rout.logger.Info(r)
	rout.router.ServeHTTP(w, r)
}

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
func (r *CustomRouter) HandleFunc(url string, f http.Handler, method string) {
	//r.router.HandlerFunc(method, url, f)
	r.router.Handle(method, url, wrapHandler(f))
}
func (rout *CustomRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rout.logger.Info(r)
	rout.router.ServeHTTP(w, r)
}
