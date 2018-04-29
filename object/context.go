package object

// Globals represents all names
type Globals struct {
	constants map[string]Object
}

// ThisContext represents `this` context
type ThisContext struct {
	globals *Globals
	vars      map[string]Object
	constants map[string]Object
}

