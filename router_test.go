package zaiuz_test

import "net/http/httptest"
import "testing"
import "io/ioutil"
import . "github.com/zaiuz/zaiuz"
import a "github.com/stretchr/testify/assert"
import "github.com/zaiuz/results"
import "github.com/zaiuz/testutil"

const TextForGet string = "GETGETGET"
const TextForPost string = "POSTPOSTPOST"
const TextForGetPost string = TextForGet + TextForPost

func TestNewRouter(t *testing.T) {
	router := NewRouter()
	a.Nil(t, router.Parent(), "new router should not have parent.")
	a.NotNil(t, router.Router(), "mux router not initialized.")
	a.Equal(t, cap(router.Filters()), InitialContextCapacity,
		"filters slice capacity mismatch.")
}

func TestIncludes(t *testing.T) {
	const UrlPath = "/subsection"

	router := NewRouter()
	a.Empty(t, router.Filters(), "filter list not empty initially.")

	subrouter := router.Subrouter(UrlPath)
	a.Empty(t, subrouter.Filters(), "empty router's subrouter filter list not empty.")

	filter := DudFilter()
	router.Include(filter)
	a.Equal(t, len(router.Filters()), 1, "module list has incorrect length.")
	a.Equal(t, len(subrouter.Filters()), 1, "subrouter module list has incorrect length.")
	a.Equal(t, router.Filters()[0], filter, "module list has wrong reference.")
	a.Equal(t, subrouter.Filters()[0], filter, "subrouter module list has wrong reference.")

	filter = DudFilter()
	subrouter.Include(filter)
	a.Equal(t, len(router.Filters()), 1, "subrouter Include() should not effects parent.")
	a.Equal(t, len(subrouter.Filters()), 2, "subrouter Include() does not take effect.")
	a.Equal(t, subrouter.Filters()[1], filter, "subrouter module list has wrong reference.")
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
		router.Static("/files", "./public")
	})
	defer server.Close()

	content, e := ioutil.ReadFile("./public/test_blob.bin")
	a.NoError(t, e)

	testutil.HttpGet(t, server.URL+"/files/test_blob.bin").Expect(200, string(content))
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

func TestModuleInvocation(t *testing.T) {
	outer, inner := testutil.NewTestFilter(), testutil.NewTestFilter()

	// TODO: Test integration with static files.
	server := newTestServer(func(router *Router) {
		router.Include(outer.Filter)
		router.Get("/", stringAction(TextForGet))

		section := router.Subrouter("/section")
		section.Include(inner.Filter)
		section.Get("/", stringAction(TextForGet))
	})
	defer server.Close()

	outer.Reset()
	inner.Reset()

	testutil.HttpGet(t, server.URL).Expect(200, TextForGet)
	a.True(t, outer.Called, "outer filter not called.")
	a.True(t, outer.Finished, "outer filter does not finish.")
	a.False(t, inner.Called, "inner filter not called.")

	outer.Reset()
	inner.Reset()

	testutil.HttpGet(t, server.URL+"/section/").Expect(200, TextForGet)
	a.True(t, outer.Called, "outer filter not called.")
	a.True(t, inner.Called, "inner filter not called.")
	a.True(t, outer.CallTime.Before(inner.CallTime), "inner filter called prematurely.")
	a.True(t, outer.Finished, "outer filter does not finish.")
	a.True(t, inner.Finished, "inner filter does not finish.")
	a.True(t, outer.FinishTime.After(inner.FinishTime), "inner filter finish prematurely.")
}

func newTestServer(setup func(router *Router)) *httptest.Server {
	router := NewRouter()
	server := httptest.NewServer(router)
	setup(router)
	return server
}

func stringAction(text string) Action {
	result := results.NewStringResult(200, text)
	return Action(func(c *Context) Result {
		return result
	})
}
