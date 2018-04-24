package eval

import "github.com/pmukhin/gophp/object"

// returnObject is a wrapper for returned objects
// to make easier to recognize when function
// execution is over
type returnObject struct {
	value object.Object
}

// Class ...
func (returnObject) Class() object.Class { panic("this function should not ever be called") }

// Id ...
func (returnObject) Id() string { panic("not implemented") }
