package zaiuz

import "fmt"
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

func NewRouter() *Router {
	modules := make([]Module, 0, InitialContextCapacity)
	return &Router{nil, mux.NewRouter(), modules}
}

// Parent() method returns the parent router if this is a child router, or nil otherwise.
func (router *Router) Parent() *Router {
	return router.parent
}

// Router() method returns the internal mux.Router from gorilla web toolkit for direct
// access. This method is not recommended unless you want to access specifici
// functionality provided by the gorilla web toolkit that does not yet have an equivalent
// in zaiuz.
func (router *Router) Router() *mux.Router {
	return router.router
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

func (router *Router) Static(urlPath, filePath string) *Router {
	// TODO: Allow action shim to still be applied
	fs := http.StripPrefix(urlPath, http.FileServer(http.Dir(filePath)))
	router.router.PathPrefix(urlPath).Methods("GET").Handler(fs)
	return router
}

func (router *Router) actionShim(action Action) func(http.ResponseWriter, *http.Request) {
	modules := router.Modules() // resolve module list w/ parents immediately

	return func(w http.ResponseWriter, r *http.Request) {
		ctx := NewContext(w, r)
		defer func() {
			if r := recover(); r != nil {
				// TODO: Errors after WriteHead?
				w.WriteHeader(500)
				w.Write([]byte(fmt.Sprintf("%s", r)))
			}
		}()

		for _, mod := range modules {
			// TODO: Detach should be given chance to recover from errors.
			e := mod.Attach(ctx)
			if e != nil {
				panic(e) // TODO: Better handover to http pkg??
			}

			defer func(m Module) {
				e := m.Detach(ctx)
				if e != nil {
					panic(e) // TODO: Better handover to http pkg?
				}
			}(mod)
		}

		action(ctx)
	}
}
