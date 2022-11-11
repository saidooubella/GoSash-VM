package object

type StringObject struct {
	Value string
}

func NewUTF8(value []uint8) Object {
	return &StringObject{string(value)}
}

func NewString(value string) Object {
	return &StringObject{value}
}

func (object *StringObject) Equals(other *StringObject) *BoolObject {
	return FromBool(object.Value == other.Value)
}

func (object *StringObject) String() string {
	return object.Value
}

func (object *StringObject) Type() ObjectType {
	return StringType
}
