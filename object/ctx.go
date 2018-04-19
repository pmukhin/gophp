package object

import (
	"fmt"
)

type Context interface {
	Outer() Context

	GetContextVar(string) (Object, error)
	SetContextVar(string, Object)
	GetFunctionTable() *FunctionTable
}

type context struct {
	outer         *context
	scope         map[string]Object
	functionTable *FunctionTable
}

func (c context) Outer() Context {
	return c.outer
}

func (c context) GetContextVar(name string) (Object, error) {
	if v, ok := c.scope[name]; ok {
		return v, nil
	}
	return nil, fmt.Errorf("name '%s' is not defined", name)
}

func (c *context) SetContextVar(name string, value Object) {
	c.scope[name] = value
}

func (c context) GetFunctionTable() *FunctionTable {
	return c.functionTable
}

func NewContext(outer Context, table *FunctionTable) Context {
	c := new(context)
	if o, ok := outer.(*context); ok {
		c.outer = o
	}
	c.scope = make(map[string]Object)
	// init function table
	c.functionTable = table

	return c
}
