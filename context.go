package zaiuz

import "net/http"

// Initial capacity for the context internal storage.
const InitialContextCapacity = 4

// The basic Context structure encapsulates request and response objects that would
// normally be supplied to a http.ServeHTTP method as well as encapsulating a small map
// for storing and passing data between modules inside a single request execution context.
type Context struct {
	objects        map[string]interface{}
	Request        *http.Request
	ResponseWriter http.ResponseWriter
}

// Creates a new context.
func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{nil, r, w}
}

// Gets an object from the context tied to the specified key.
func (c *Context) Get(key string) interface{} {
	if c.objects == nil {
		return nil
	} else {
		return c.objects[key]
	}
}

// Gets an object from the context tied to the specified key, also returns an additional
// boolean indicating wether the get was successful.
func (c *Context) GetOk(key string) (interface{}, bool) {
	if c.objects == nil {
		return nil, false
	} else {
		result, ok := c.objects[key]
		return result, ok
	}
}

// Checks wether or not the context contains the specified key. You can also use GetOk to
// perform the checking and retreival in one operation.
func (c *Context) Has(key string) bool {
	_, ok := c.objects[key]
	return ok
}

// Saves an object the the context and associate it to the specified key.
func (c *Context) Set(key string, value interface{}) {
	if c.objects == nil {
		c.objects = make(map[string]interface{}, InitialContextCapacity)
	}
	c.objects[key] = value
}

// Removes an object tied to the specified key from the context.
func (c *Context) Delete(key string) {
	delete(c.objects, key)
}
