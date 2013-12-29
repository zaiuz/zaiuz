package zaiuz

type Module interface {
	Attach(ctx *Context) error
	Detach(ctx *Context) error
	// TODO: Recover(ctx *context.Context, e error)
}
