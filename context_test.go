package zaiuz_test

import "testing"
import "github.com/zaiuz/testutil"
import "github.com/zaiuz/results"
import . "github.com/zaiuz/zaiuz"
import a "github.com/stretchr/testify/assert"

func TestNewContext(t *testing.T) {
	response, request := testutil.NewTestRequestPair()
	result := NewContext(response, request)

	a.Equal(t, result.Request, request, "request not saved.")
	a.Equal(t, result.ResponseWriter, response, "response writer not saved.")
}

func TestContextObjects(t *testing.T) {
	response, request := testutil.NewTestRequestPair()
	context := NewContext(response, request)

	const KEY = "context_item"
	type Dummy struct{}

	has := context.Has(KEY)
	result := context.Get(KEY)

	a.False(t, has, "Has(key) should be false.")
	a.Nil(t, result, "result value should be nil.")

	result, ok := context.GetOk(KEY)
	a.False(t, ok, "ok true for non-existent key.")
	a.Nil(t, result, "result value should be nil.")

	value := &Dummy{}
	context.Set(KEY, value)
	has = context.Has(KEY)
	result, ok = context.Get(KEY).(*Dummy)

	a.True(t, has, "Has(key) should be true.")
	a.True(t, ok, "type assertion failed.")
	a.IsType(t, result, &Dummy{}, "return value is of wrong type.")
	a.Equal(t, result, value, "return value is a different instance.")

	result, ok = context.GetOk(KEY)
	a.True(t, ok, "ok false for key that definitely exists.")
	a.IsType(t, result, &Dummy{}, "return value is of wrong type.")
	a.Equal(t, result, value, "return value is a different instance.")

	context.Delete(KEY)
	has = context.Has(KEY)
	result = context.Get(KEY)

	a.False(t, has, "Has(key) should be false.")
	a.Nil(t, result, "return value should be nil.")

	result, ok = context.GetOk(KEY)
	a.False(t, ok, "ok true for deleted key.")
	a.Nil(t, result, "result non-nil for deleted key.")
}

func TestNewContext_WithRouteVars(t *testing.T) {
	c := make(chan *Context)
	done := make(chan bool)

	// TODO: How to test this without relying on the mux router or the test server?
	server := newTestServer(func(router *Router) {
		router.Get("/with/{id}", func(ctx *Context) Result {
			c <- ctx
			done <- true
			return results.String(200, "OK")
		})
	})

	go func() {
		testutil.HttpGet(t, server.URL+"/with/asdf").Expect(200, "OK")
		done <- true
	}()

	ctx := <-c
	raw, ok := ctx.GetOk("id")
	a.True(t, ok, "router params not saved to context.")

	result, ok := raw.(string)
	a.True(t, ok, "router params not a string.")
	a.Equal(t, result, "asdf", "saved router params has wrong value.")

	<-done
	<-done
}
