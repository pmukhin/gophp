package object

var (
	arrayClass = internalClass{
		name:      "Array",
		final:     false,
		abstract:  false,
		methodSet: newMethodSet(nil),
	}
)

// NewArray ...
func NewArray(os ...Object) (Object, error) {
	array := new(ArrayObject)
	array.Values = make([]Object, len(os), 32)

	for i, ob := range os {
		array.Values[i] = ob
	}

	return array, nil
}

type ArrayObject struct {
	Values []Object
}

func (ArrayObject) Class() Class {
	return arrayClass
}

func (ArrayObject) Id() string {
	panic("implement me")
}
