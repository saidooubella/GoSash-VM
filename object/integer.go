package object

import "strconv"

type IntObject struct {
	Value int32
}

func NewInt(value int32) Object {
	return &IntObject{value}
}

func (object *IntObject) Less(other *IntObject) *BoolObject {
	return FromBool(object.Value < other.Value)
}

func (object *IntObject) Greater(other *IntObject) *BoolObject {
	return FromBool(object.Value > other.Value)
}

func (object *IntObject) Equals(other *IntObject) *BoolObject {
	return FromBool(object.Value == other.Value)
}

func (object *IntObject) Plus(other *IntObject) Object {
	return &IntObject{object.Value + other.Value}
}

func (object *IntObject) Sub(other *IntObject) Object {
	return &IntObject{object.Value - other.Value}
}

func (object *IntObject) Div(other *IntObject) Object {
	return &IntObject{object.Value / other.Value}
}

func (object *IntObject) Mul(other *IntObject) Object {
	return &IntObject{object.Value * other.Value}
}

func (object *IntObject) Neg() Object {
	return &IntObject{-object.Value}
}

func (object *IntObject) String() string {
	return strconv.FormatInt(int64(object.Value), 10)
}

func (object *IntObject) Type() ObjectType {
	return IntType
}
