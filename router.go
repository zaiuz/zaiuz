package zaiuz

import "net/http"
import "github.com/gorilla/mux"

type Router struct {
	parent  *Router
	router  *mux.Router
	modules []Module
}

var _ http.Handler = &Router{}

type Action func(ctx *Context)

func NewRouter() *Router {
	modules := make([]Module, 0, InitialContextCapacity)
	return &Router{nil, mux.NewRouter(), modules}
}

func (router *Router) Modules() []Module {
	if router.parent == nil {
		return router.modules
	} else {
		return append(router.parent.modules, router.modules...)
	}
}

func (router *Router) Subrouter(path string) *Router {
	subrouter := router.router.PathPrefix(path).Subrouter()
	result := &Router{router, subrouter, []Module{}}
	return result
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router.router.ServeHTTP(w, r)
}

func (router *Router) Include(m ...Module) {
	router.modules = append(router.modules, m...)
}

func (router *Router) Get(path string, action Action) *Router {
	handler := router.actionShim(action)
	router.router.Path(path).Methods("GET").HandlerFunc(handler)
	return router
}

func (router *Router) Post(path string, action Action) *Router {
	handler := router.actionShim(action)
	router.router.Path(path).Methods("POST").HandlerFunc(handler)
	return router
}

func (router *Router) GetPost(path string, action Action) *Router {
	handler := router.actionShim(action)
	router.router.Path(path).Methods("GET", "POST").HandlerFunc(handler)
	return router
}

func (router *Router) actionShim(action Action) func(http.ResponseWriter, *http.Request) {
	modules := router.Modules() // resolve module list w/ parents

	return func(w http.ResponseWriter, r *http.Request) {
		ctx := NewContext(w, r)

		for _, mod := range modules {
			// TODO: Handle attach/detach errors. panic?
			// TODO: Detach should be given chance to recover from errors.
			mod.Attach(ctx)
			defer mod.Detach(ctx)
		}

		action(ctx)
	}
}
