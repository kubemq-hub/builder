package main

type Context struct {
}

func (c *Context) Set() error {
	return nil
}

func (c *Context) Get() string {
	return "local_file"
}
