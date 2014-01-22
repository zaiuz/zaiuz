package zaiuz_test

import "net/http/httptest"
import "testing"
import "io/ioutil"
import "fmt"
import "time"
import . "github.com/zaiuz/zaiuz"
import a "github.com/stretchr/testify/assert"
import "github.com/zaiuz/testutil"

const TextForGet string = "GETGETGET"
const TextForPost string = "POSTPOSTPOST"
const TextForGetPost string = TextForGet + TextForPost

type DummyModule struct {
	attachCalled bool
	attachTime   time.Time
	detachCalled bool
	detachTime   time.Time
}

type PanicModule struct{}

var _ Module = new(DummyModule)
var _ Module = new(PanicModule)

func (m *DummyModule) Reset() {
	zero := time.Unix(0, 0)
	m.attachCalled, m.detachCalled = false, false
	m.attachTime, m.detachTime = zero, zero
}

func (m *DummyModule) Attach(c *Context) error {
	m.attachCalled = true
	m.attachTime = time.Now()
	return nil
}

func (m *DummyModule) Detach(c *Context) error {
	m.detachCalled = true
	m.detachTime = time.Now()
	return nil
}

func (p *PanicModule) Attach(c *Context) error {
	return fmt.Errorf("PanicModule test error.")
}

func (p *PanicModule) Detach(c *Context) error {
	return fmt.Errorf("PanicModule test error.")
}

func TestNewRouter(t *testing.T) {
	router := NewRouter()
	a.Nil(t, router.Parent(), "new router should not have parent.")
	a.NotNil(t, router.Router(), "mux router not initialized.")
	a.Equal(t, cap(router.Modules()), InitialContextCapacity,
		"modules slice capacity mismatch.")
}

func TestIncludes(t *testing.T) {
	const UrlPath = "/subsection"

	router := NewRouter()
	a.Empty(t, router.Modules(), "module list not empty initially.")

	subrouter := router.Subrouter(UrlPath)
	a.Empty(t, subrouter.Modules(), "empty router's subrouter module list not empty.")

	mod := new(DummyModule)
	router.Include(mod)
	a.Equal(t, len(router.Modules()), 1, "module list has incorrect length.")
	a.Equal(t, len(subrouter.Modules()), 1, "subrouter module list has incorrect length.")
	a.Equal(t, router.Modules()[0], mod, "module list has wrong reference.")
	a.Equal(t, subrouter.Modules()[0], mod, "subrouter module list has wrong reference.")

	mod = new(DummyModule)
	subrouter.Include(mod)
	a.Equal(t, len(router.Modules()), 1, "subrouter Include() should not effects parent.")
	a.Equal(t, len(subrouter.Modules()), 2, "subrouter Include() does not take effect.")
	a.Equal(t, subrouter.Modules()[1], mod, "subrouter module list has wrong reference.")
}

func TestSubrouter(t *testing.T) {
	router := NewRouter()
	subrouter := router.Subrouter("/section")

	a.NotNil(t, subrouter, "subrouter creation failed.")
	a.NotEqual(t, router, subrouter, "subrouter must not equals parent.")
	a.NotNil(t, subrouter.Parent(), "subrouter should have parent.")
	a.Equal(t, subrouter.Parent(), router, "subrouter has wrong parent.")
}

func TestSimpleGet(t *testing.T) {
	server := newTestServer(func(router *Router) {
		router.Get("/", stringAction(TextForGet))
	})
	defer server.Close()

	testutil.HttpGet(t, server.URL).Expect(200, TextForGet)
}

func TestMethodRouting(t *testing.T) {
	server := newTestServer(func(router *Router) {
		router.Get("/", stringAction(TextForGet))
		router.Post("/", stringAction(TextForPost))
		router.GetPost("/twins", stringAction(TextForGetPost))
	})
	defer server.Close()

	testutil.HttpGet(t, server.URL).Expect(200, TextForGet)
	testutil.HttpPost(t, server.URL, nil).Expect(200, TextForPost)
	testutil.HttpGet(t, server.URL+"/twins").Expect(200, TextForGetPost)
	testutil.HttpPost(t, server.URL+"/twins", nil).Expect(200, TextForGetPost)
}

func TestStaticFiles(t *testing.T) {
	server := newTestServer(func(router *Router) {
		router.Static("/files", "./testviews")
	})
	defer server.Close()

	content, e := ioutil.ReadFile("./testviews/single.html")
	a.NoError(t, e)

	testutil.HttpGet(t, server.URL+"/files/single.html").Expect(200, string(content))
}

func TestSubrouterRouting(t *testing.T) {
	server := newTestServer(func(router *Router) {
		router.Get("/", stringAction(TextForGet))
		section := router.Subrouter("/section")
		section.Get("/", stringAction(TextForGet+"inner"))
	})
	defer server.Close()

	// TODO: Handle missing '/'
	testutil.HttpGet(t, server.URL).Expect(200, TextForGet)
	testutil.HttpGet(t, server.URL+"/section/").Expect(200, TextForGet+"inner")
}

func TestModulePanic(t *testing.T) {
	server := newTestServer(func(router *Router) {
		router.Include(new(PanicModule))
		router.Get("/", stringAction(TextForGet))
	})

	testutil.HttpGet(t, server.URL).ExpectPattern(500, "^PanicModule")
}

func TestModuleInvocation(t *testing.T) {
	outer, inner := new(DummyModule), new(DummyModule)

	// TODO: Test that static files still cause modules  to be invoked.
	server := newTestServer(func(router *Router) {
		router.Include(outer)
		router.Get("/", stringAction(TextForGet))

		section := router.Subrouter("/section")
		section.Include(inner)
		section.Get("/", stringAction(TextForGet))
	})
	defer server.Close()

	outer.Reset()
	inner.Reset()

	testutil.HttpGet(t, server.URL).Expect(200, TextForGet)
	a.True(t, outer.attachCalled, "outer module not attached.")
	a.True(t, outer.detachCalled, "outer module not detached.")
	a.True(t, outer.attachTime.Before(outer.detachTime), "detached before attaching.")
	a.False(t, inner.attachCalled, "inner module incorrectly attached.")
	a.False(t, inner.detachCalled, "inner module incorrectly detached.")

	outer.Reset()
	inner.Reset()

	testutil.HttpGet(t, server.URL+"/section/").Expect(200, TextForGet)
	a.True(t, outer.attachCalled, "outer module not attached.")
	a.True(t, inner.attachCalled, "inner module not attached.")
	a.True(t, outer.attachTime.Before(inner.attachTime), "inner module attach prematurely.")
	a.True(t, outer.detachCalled, "outer module not detached.")
	a.True(t, inner.detachCalled, "inner module not detached.")
	a.True(t, outer.detachTime.After(inner.detachTime), "outer module detach prematurely.")
}

func newTestServer(setup func(router *Router)) *httptest.Server {
	router := NewRouter()
	server := httptest.NewServer(router)
	setup(router)
	return server
}

func stringAction(text string) Action {
	return Action(func(c *Context) {
		c.ResponseWriter.Write([]byte(text))
	})
}

func moduleListEquals(a, b []Module) bool {
	switch {
	case len(a) != len(b):
		return false
	case cap(a) != cap(b):
		return false
	}

	for _, ma := range a {
		for _, mb := range b {
			if ma != mb {
				return false
			}
		}
	}

	return true
}
