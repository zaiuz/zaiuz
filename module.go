package zaiuz

// Modules, or middlewares, provides a way to add functionality that is orthogonal the
// normal request processing pipeline such as session or database transaction which
// requires saving/closing at the end of each request.
//
// Functionality that effect multiple actions at once or that which requires a "close"
// action at the end are best implemented as Modules to reduce the number of code
// duplication.
type Module interface {
	Attach(ctx *Context) error
	Detach(ctx *Context) error
	Recover(ctx *Context, e interface{}) bool
	// TODO: Recover(ctx *context.Context, e error)
}
