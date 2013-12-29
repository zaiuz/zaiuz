package zaiuz

import "net/http"

const InitialContextCapacity = 4

type Context struct {
	objects        map[string]interface{}
	Request        *http.Request
	ResponseWriter http.ResponseWriter
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{nil, r, w}
}

func (c *Context) Get(key string) interface{} {
	if c.objects == nil {
		return nil
	} else {
		return c.objects[key]
	}
}

func (c *Context) GetOk(key string) (interface{}, bool) {
	if c.objects == nil {
		return nil, false
	} else {
		result, ok := c.objects[key]
		return result, ok
	}
}

func (c *Context) Has(key string) bool {
	_, ok := c.objects[key]
	return ok
}

func (c *Context) Set(key string, value interface{}) {
	if c.objects == nil {
		c.objects = make(map[string]interface{}, InitialContextCapacity)
	}
	c.objects[key] = value
}

func (c *Context) Delete(key string) {
	delete(c.objects, key)
}
