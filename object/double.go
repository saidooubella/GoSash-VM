package object

import "strconv"

type DoubleObject struct {
	Value float64
}

func NewDouble(value float64) Object {
	return &DoubleObject{value}
}

func (object *DoubleObject) Less(other *DoubleObject) *BoolObject {
	return FromBool(object.Value < other.Value)
}

func (object *DoubleObject) Greater(other *DoubleObject) *BoolObject {
	return FromBool(object.Value > other.Value)
}

func (object *DoubleObject) Equals(other *DoubleObject) *BoolObject {
	return FromBool(object.Value == other.Value)
}

func (object *DoubleObject) Plus(other *DoubleObject) Object {
	return &DoubleObject{object.Value + other.Value}
}

func (object *DoubleObject) Sub(other *DoubleObject) Object {
	return &DoubleObject{object.Value - other.Value}
}

func (object *DoubleObject) Div(other *DoubleObject) Object {
	return &DoubleObject{object.Value / other.Value}
}

func (object *DoubleObject) Mul(other *DoubleObject) Object {
	return &DoubleObject{object.Value * other.Value}
}

func (object *DoubleObject) Neg() Object {
	return &DoubleObject{-object.Value}
}

func (object *DoubleObject) String() string {
	return strconv.FormatFloat(object.Value, 'f', -1, 64)
}

func (object *DoubleObject) Type() ObjectType {
	return DoubleType
}
