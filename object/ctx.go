package object

import (
	"fmt"
)

type Context interface {
	Outer() Context

	// globs
	SetGlobal(string, Object) error
	GetGlobal(name string) (Object, error)

	// vars
	GetContextVar(string) (Object, error)
	SetContextVar(string, Object) error

	GetFunctionTable() *FunctionTable
}

type context struct {
	outer         *context
	scope         map[string]Object
	globalsTable  map[string]Object
	functionTable *FunctionTable
}

func (c context) SetGlobal(name string, value Object) error {
	if _, ok := c.globalsTable[name]; ok {
		return fmt.Errorf("can not redeclare const '%s'", name)
	}
	c.globalsTable[name] = value
	return nil
}

func (c context) GetGlobal(name string) (Object, error) {
	if v, ok := c.globalsTable[name]; ok {
		return v, nil
	}
	return nil, fmt.Errorf("name '%s' is not defined", name)
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

func (c *context) SetContextVar(name string, value Object) error {
	c.scope[name] = value
	return nil
}

// GetFunctionTable ...
func (c context) GetFunctionTable() *FunctionTable {
	return c.functionTable
}

func NewContext(outer Context, table *FunctionTable) Context {
	c := new(context)
	if o, ok := outer.(*context); ok {
		c.outer = o
	}
	c.scope = make(map[string]Object)
	c.globalsTable = make(map[string]Object)
	// init function table
	c.functionTable = table

	return c
}
