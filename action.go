package zaiuz

// Action is the main interaction unit of zaiuz. This is zaiuz's analog to the standard
// http.ServeHTTP method. Most methods that work on the Context should follows the same
// function signature.
type Action func(ctx *Context) // TODO: Return Result
