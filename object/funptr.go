package object

import "fmt"

type FunPtrObject struct {
	Ptr   uint16
	Arity uint8
}

func NewFunPtr(ptr uint16, arity uint8) Object {
	return &FunPtrObject{ptr, arity}
}

func (object *FunPtrObject) String() string {
	return fmt.Sprintf("FunPtr(ptr: %d, arity: %d)", object.Ptr, object.Arity)
}

func (object *FunPtrObject) Type() ObjectType {
	return FunPtrType
}
