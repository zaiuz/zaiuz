package zaiuz

type Result interface{
	Execute(c *Context) error
}

type funcResult struct{
	execute func(c *Context) error
}

func (result *funcResult) Execute(c *Context) error {
	return result.execute(c)
}

func ResultFunc(execute func(c *Context) error) Result {
	return &funcResult{execute}
}

func DudResult() Result {
	return ResultFunc(func(c *Context) error {
		return nil
	})
}

