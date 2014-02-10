package zaiuz

// ActionFilter provides a way to wraps an existing actions to modifies or adds to the
// functionality. ActionFilter is functionally equivalent to traditional middlewares.
type Filter func(action Action) Action

// DudFilter returns a Filter that have no effect. Mostly only useful for testing.
func DudFilter() Filter {
	return func(action Action) Action { return action }
}
