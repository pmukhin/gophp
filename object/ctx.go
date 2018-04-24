package object

import (
	"fmt"
)

type localStorage struct {
	vars   map[string]Object
	parent *localStorage
}

func (l *localStorage) isSet(name string) bool {
	_, ok := l.vars[name]
	return ok
}

func (l *localStorage) Set(name string, o Object) error {
	// first let'par try to found this name somewhere in the outer context
	par := l
	for par != nil && !par.isSet(name) {
		par = par.parent
	}
	// if it'par there par would not be null
	if par == nil {
		l.vars[name] = o
	} else {
		par.vars[name] = o
	}

	return nil
}

func (l *localStorage) Get(name string) (Object, error) {
	s := l
	for s != nil && !s.isSet(name) {
		s = s.parent
	}
	if s != nil {
		return s.vars[name], nil
	}
	return Null, nil
}

func (l *localStorage) SetParent(s *localStorage) {
	l.parent = s
}

func newLocalStorage() *localStorage {
	return &localStorage{vars: make(map[string]Object)}
}

type Context interface {
	// globs
	SetGlobal(string, Object) error
	GetGlobal(name string) (Object, error)

	// vars
	GetContextVar(string) (Object, error)
	SetContextVar(string, Object) error
	Scope() *localStorage
}

type context struct {
	scope        *localStorage
	globalsTable map[string]Object
}

func (c *context) Scope() *localStorage {
	return c.scope
}

func (c *context) SetGlobal(name string, value Object) error {
	if _, ok := c.globalsTable[name]; ok {
		return fmt.Errorf("can not redeclare const '%s'", name)
	}
	c.globalsTable[name] = value
	return nil
}

func (c *context) GetGlobal(name string) (Object, error) {
	if v, ok := c.globalsTable[name]; ok {
		return v, nil
	}
	return nil, fmt.Errorf("name '%s' is not defined", name)
}

func (c *context) GetContextVar(name string) (Object, error) {
	return c.scope.Get(name)
}

func (c *context) SetContextVar(name string, value Object) error {
	return c.scope.Set(name, value)
}

func CloneContext(ctx Context, parentScope *localStorage) Context {
	c := new(context)
	c.scope = newLocalStorage()
	c.scope.SetParent(parentScope)

	// copy all global names
	c.globalsTable = ctx.(*context).globalsTable

	return c
}

func NewContext(parent *localStorage) Context {
	c := new(context)
	// init scope
	c.scope = newLocalStorage()
	c.scope.SetParent(parent)

	c.globalsTable = make(map[string]Object)

	return c
}
