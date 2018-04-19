package object

var (
	BooleanClass = internalClass{
		name:      "Boolean",
		final:     true,
		abstract:  false,
		methodSet: newMethodSet(nil),
	}

	True  = &BooleanObject{Value: true}
	False = &BooleanObject{Value: false}
)

type BooleanObject struct {
	Value bool
}

func (BooleanObject) Class() Class {
	return BooleanClass
}

func (BooleanObject) Id() string {
	panic("implement me")
}
