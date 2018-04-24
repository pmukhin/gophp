package object

var (
	nullMethods = map[string]Method{
		"__toString": newMethod(func(this Object, args ...Object) (Object, error) {
			return nil, nil
		}, VisibilityPublic),
	}

	classNull = &InternalClass{
		name:      "Null",
		final:     true,
		methodSet: newMethodSet(nullMethods),
	}

	Null = &NullObject{}
)

type NullObject struct{}

func (NullObject) Class() Class { return classNull }

func (NullObject) Id() string { panic("id") }
