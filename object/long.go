package object

import "strconv"

type LongObject struct {
	Value int64
}

func NewLong(value int64) Object {
	return &LongObject{value}
}

func (object *LongObject) Less(other *LongObject) *BoolObject {
	return FromBool(object.Value < other.Value)
}

func (object *LongObject) Greater(other *LongObject) *BoolObject {
	return FromBool(object.Value > other.Value)
}

func (object *LongObject) Equals(other *LongObject) *BoolObject {
	return FromBool(object.Value == other.Value)
}

func (object *LongObject) Plus(other *LongObject) Object {
	return &LongObject{object.Value + other.Value}
}

func (object *LongObject) Sub(other *LongObject) Object {
	return &LongObject{object.Value - other.Value}
}

func (object *LongObject) Div(other *LongObject) Object {
	return &LongObject{object.Value / other.Value}
}

func (object *LongObject) Mul(other *LongObject) Object {
	return &LongObject{object.Value * other.Value}
}

func (object *LongObject) Neg() Object {
	return &LongObject{-object.Value}
}

func (object *LongObject) String() string {
	return strconv.FormatInt(object.Value, 10)
}

func (object *LongObject) Type() ObjectType {
	return LongType
}
