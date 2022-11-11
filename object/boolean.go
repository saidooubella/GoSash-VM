package object

import "strconv"

var (
	BoolFalse = &BoolObject{false}
	BoolTrue  = &BoolObject{true}
)

type BoolObject struct {
	Value bool
}

func FromBool(value bool) *BoolObject {
	if value {
		return BoolTrue
	}
	return BoolFalse
}

func (object *BoolObject) Equals(other *BoolObject) *BoolObject {
	return FromBool(object.Value == other.Value)
}

func (object *BoolObject) Inv() Object {
	return FromBool(!object.Value)
}

func (object *BoolObject) String() string {
	return strconv.FormatBool(object.Value)
}

func (object *BoolObject) Type() ObjectType {
	return BoolType
}
