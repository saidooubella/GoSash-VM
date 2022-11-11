package object

import "fmt"

type ClosureObject struct {
	Func       *FunPtrObject
	FreeValues []*FreeValue
}

type FreeValue struct {
	Pointer *Object
	Value   Object
}

func NewClosure(fun *FunPtrObject, freeCount uint16) *ClosureObject {
	values := make([]*FreeValue, freeCount)
	return &ClosureObject{fun, values}
}

func (object *ClosureObject) String() string {
	return fmt.Sprintf("Closure(fun: %v, freeValues: %v)", object.Func, object.FreeValues)
}

func (object *ClosureObject) Type() ObjectType {
	return ClosureType
}
