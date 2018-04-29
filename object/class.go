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

func NewUserClass(final bool, abstract bool, methodSet MethodSet, propertySet map[string]Object) Class {
	return &UserClass{
		final:           final,
		abstract:        abstract,
		methodSet:       methodSet,
		staticMethodSet: nil,
		propertySet:     propertySet,
	}
}

type UserClass struct {
	final           bool
	abstract        bool
	methodSet       MethodSet
	staticMethodSet MethodSet
	propertySet     map[string]Object
}

func (UserClass) Name() string {
	panic("implement me")
}

func (UserClass) Constructor() Method {
	panic("implement me")
}

func (UserClass) SuperClass() Class {
	panic("implement me")
}

func (UserClass) IsFinal() bool {
	panic("implement me")
}

func (UserClass) IsAbstract() bool {
	panic("implement me")
}

func (UserClass) Methods() MethodSet {
	panic("implement me")
}

func (UserClass) StaticMethods() MethodSet {
	panic("implement me")
}

type InternalClass struct {
	name                string
	final               bool
	abstract            bool
	constructor         Method
	internalConstructor InternalConstructor
	methodSet           MethodSet
}

func (c InternalClass) InternalConstructor(value interface{}) (Object, error) {
	return c.internalConstructor(value)
}

func (c InternalClass) Class() Class {
	panic("implement me")
}

func (c InternalClass) Id() string {
	panic("implement me")
}

func (c InternalClass) Name() string {
	return c.name
}

func (c InternalClass) Constructor() Method {
	return c.constructor
}

func (InternalClass) SuperClass() Class {
	panic("implement me")
}

func (c InternalClass) IsFinal() bool {
	return c.final
}

func (c InternalClass) IsAbstract() bool {
	return c.abstract
}

func (c InternalClass) Methods() MethodSet {
	return c.methodSet
}

func (InternalClass) StaticMethods() MethodSet {
	panic("implement me")
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
