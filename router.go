package zaiuz

import "net/http"
import "github.com/gorilla/mux"

// Router is a wrapper over Gorrila web toolkit's mux router
// (http://www.gorillatoolkit.org/pkg/mux) that provides helpers that work with standard
// zaiuz Action function signature as well as Filters.
type Router struct {
	parent  *Router
	router  *mux.Router
	filters []Filter
}

var _ http.Handler = &Router{}

func NewRouter() *Router {
	filters := make([]Filter, 0, InitialContextCapacity)
	return &Router{nil, mux.NewRouter(), filters}
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
func (router *Router) Filters() []Filter {
	if router.parent == nil {
		return router.filters
	} else {
		return append(router.parent.filters, router.filters...)
	}
}

// Creates a child router. Analogous to calling mux's Router.Subrouter function but also
// carry over all the Filters included so far as well. Filters added in the subrouter only
// run inside the subrouter.
func (router *Router) Subrouter(path string) *Router {
	subrouter := router.router.PathPrefix(path).Subrouter()
	result := &Router{router, subrouter, []Filter{}}
	return result
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router.router.ServeHTTP(w, r)
}

// Include a Filter in the current router for all mapped (and future) actions.
func (router *Router) Include(filters ...Filter) {
	router.filters = append(router.filters, filters...)
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
	// apply innermost filter first
	filters := router.Filters()
	for i := range filters {
		// TODO: Detect nil
		action = filters[len(filters)-1-i](action)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		context := NewContext(w, r)
		result := action(context)
		result.Render(context)
	}
}
