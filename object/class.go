package object

type Visibility uint8

const (
	VisibilityPublic    Visibility = iota
	VisibilityProtected
	VisibilityPrivate
)

type Class interface {
	Name() string
	Constructor() Method
	SuperClass() Class
	IsFinal() bool
	IsAbstract() bool
	Methods() MethodSet
	StaticMethods() MethodSet
}

type InternalConstructor func(value interface{}) (Object, error)

type InternalClass interface {
	Class
	InternalConstructor(value interface{}) (Object, error)
}

type Method interface {
	Call(this Object, args ...Object) (Object, error)
	Visibility() Visibility
}

type StaticMethod interface {
	Call(args ...Object) (Object, error)
	Visibility() Visibility
}

type MethodSet interface {
	Find(string) Method
	All() []Method
}

type Object interface {
	Class() Class
	Id() string
}

type internalClass struct {
	name                string
	final               bool
	abstract            bool
	constructor         Method
	internalConstructor InternalConstructor
	methodSet           MethodSet
}

func (c internalClass) InternalConstructor(value interface{}) (Object, error) {
	return c.internalConstructor(value)
}

func (c internalClass) Class() Class {
	panic("implement me")
}

func (c internalClass) Id() string {
	panic("implement me")
}

func (c internalClass) Name() string {
	return c.name
}

func (c internalClass) Constructor() Method {
	return c.constructor
}

func (internalClass) SuperClass() Class {
	panic("implement me")
}

func (c internalClass) IsFinal() bool {
	return c.final
}

func (c internalClass) IsAbstract() bool {
	return c.abstract
}

func (c internalClass) Methods() MethodSet {
	return c.methodSet
}

func (internalClass) StaticMethods() MethodSet {
	panic("implement me")
}

func newInternalClass(name string, final, abstract bool, constructor Method, ic InternalConstructor) InternalClass {
	return &internalClass{
		name:                name,
		final:               final,
		abstract:            abstract,
		constructor:         constructor,
		internalConstructor: ic,
	}
}

type methodSet struct {
	nameMap map[string]Method
}

func (ms methodSet) Find(name string) Method {
	return ms.nameMap[name]
}

func (methodSet) All() []Method {
	panic("implement me")
}

func newMethodSet(nameMap map[string]Method) MethodSet {
	return &methodSet{nameMap: nameMap}
}

type method struct {
	name string
	f    func(this Object, args ...Object) (Object, error)
	vis  Visibility
}

func (m method) Call(this Object, args ...Object) (Object, error) {
	return m.f(this, args...)
}

func (m method) Visibility() Visibility {
	return m.vis
}

func newMethod(f func(this Object, args ...Object) (Object, error), vis Visibility) Method {
	return &method{f: f, vis: vis}
}
