package zaiuz

import "net/http"
import "github.com/gorilla/mux"

// Router is a thin wrapper over Gorrila web toolkit's mux router
// (http://www.gorillatoolkit.org/pkg/mux) that provides helper for common routing actions
// such as Get/Post actions. Zaiuz's router also requires the Action function signature.
//
// Zaiuz's router also provide mechanism for executing Module Attach/Detach and other
// related hooks without having to explicitly call them from your Action code.
type Router struct {
	parent  *Router
	router  *mux.Router
	modules []Module
}

var _ http.Handler = &Router{}

// Action is the main interaction unit of zaiuz. This is zaiuz's analog to the standard
// http.ServeHTTP method. Most methods that work on the Context should follows the same
// function signature.
type Action func(ctx *Context)

func NewRouter() *Router {
	modules := make([]Module, 0, InitialContextCapacity)
	return &Router{nil, mux.NewRouter(), modules}
}

// Retreive all modules included into this router so far. Also resolve parent's list of
// modules if called from a subrouter.
func (router *Router) Modules() []Module {
	if router.parent == nil {
		return router.modules
	} else {
		return append(router.parent.modules, router.modules...)
	}
}

// Creates a child router. Analogous to calling mux's Router.Subrouter function.
// Additionally, any zaiuz's Module added to a subrouter are scoped to only that subrouter
// and any child subrouters.
func (router *Router) Subrouter(path string) *Router {
	subrouter := router.router.PathPrefix(path).Subrouter()
	result := &Router{router, subrouter, []Module{}}
	return result
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router.router.ServeHTTP(w, r)
}

// Include a Module in the current router for all mapped (and future) actions.
func (router *Router) Include(m ...Module) {
	router.modules = append(router.modules, m...)
}

// Maps an Action to the specified URL Path and HTTP GET method.
func (router *Router) Get(path string, action Action) *Router {
	handler := router.actionShim(action)
	router.router.Path(path).Methods("GET").HandlerFunc(handler)
	return router
}

// Maps an Action to the specified URL Path and HTTP POST method.
func (router *Router) Post(path string, action Action) *Router {
	handler := router.actionShim(action)
	router.router.Path(path).Methods("POST").HandlerFunc(handler)
	return router
}

// Maps an Action to the specified URL Path and HTTP GET _and_ POST method.
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
