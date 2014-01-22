package zaiuz

type Result interface {
	Render(c *Context) error
}
