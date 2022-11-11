package vm

import (
	"demo/object"
	"fmt"
)

const initialStackSize = 256

type ObjectStack struct {
	stack []object.Object
	size  int
}

func NewObjectStack() ObjectStack {
	return ObjectStack{
		stack: make([]object.Object, initialStackSize),
		size:  0,
	}
}

func (stack *ObjectStack) Push(object object.Object) {
	stack.stack = ensureSpace(stack.stack, stack.size, 1)
	stack.stack[stack.size] = object
	stack.size += 1
}

func (stack *ObjectStack) Peek() object.Object {
	ensureNotEmpty(stack.size)
	return stack.stack[stack.size-1]
}

func (stack *ObjectStack) PeekPtr() *object.Object {
	ensureNotEmpty(stack.size)
	return &stack.stack[stack.size-1]
}

func (stack *ObjectStack) Pop() object.Object {
	ensureNotEmpty(stack.size)
	stack.size -= 1
	temp := stack.stack[stack.size]
	stack.stack[stack.size] = nil
	return temp
}

func (stack *ObjectStack) Get(index int) object.Object {
	checkIndexBounds(stack.size, index)
	return stack.stack[index]
}

func (stack *ObjectStack) GetPtr(index int) *object.Object {
	checkIndexBounds(stack.size, index)
	return &stack.stack[index]
}

func (stack *ObjectStack) Set(index int, object object.Object) {
	checkIndexBounds(stack.size, index)
	stack.stack[index] = object
}

func (stack *ObjectStack) Size() int {
	return stack.size
}

func (stack *ObjectStack) Reset() {
	for i := stack.size - 1; i >= 0; i-- {
		stack.stack[i] = nil
	}
	stack.size = 0
}

func (stack *ObjectStack) Debug() {
	fmt.Print("[")
	for index := 0; index < stack.size; index++ {
		fmt.Printf("'%s'", stack.stack[index].String())
		if index < stack.size-1 {
			fmt.Print(", ")
		}
	}
	fmt.Print("]\n")
}
