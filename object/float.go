package object

import "strconv"

type FloatObject struct {
	Value float32
}

func NewFloat(value float32) Object {
	return &FloatObject{value}
}

func (object *FloatObject) Less(other *FloatObject) *BoolObject {
	return FromBool(object.Value < other.Value)
}

func (object *FloatObject) Greater(other *FloatObject) *BoolObject {
	return FromBool(object.Value > other.Value)
}

func (object *FloatObject) Equals(other *FloatObject) *BoolObject {
	return FromBool(object.Value == other.Value)
}

func (object *FloatObject) Plus(other *FloatObject) Object {
	return &FloatObject{object.Value + other.Value}
}

func (object *FloatObject) Sub(other *FloatObject) Object {
	return &FloatObject{object.Value - other.Value}
}

func (object *FloatObject) Div(other *FloatObject) Object {
	return &FloatObject{object.Value / other.Value}
}

func (object *FloatObject) Mul(other *FloatObject) Object {
	return &FloatObject{object.Value * other.Value}
}

func (object *FloatObject) Neg() Object {
	return &FloatObject{-object.Value}
}

func (object *FloatObject) String() string {
	return strconv.FormatFloat(float64(object.Value), 'f', -1, 32)
}

func (object *FloatObject) Type() ObjectType {
	return FloatType
}
