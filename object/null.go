package object

var Null = &NullObject{}

type NullObject struct{}

func (NullObject) Class() Class { panic("implement me") }

func (NullObject) Id() string {
	return "0"
}
