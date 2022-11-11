package vm

import "demo/object"

const initialFrameSize = 256

type StackFrame struct {
	// It can be `nil`
	closure *object.ClosureObject
	bp      int
	ip      int
}

type CallStack struct {
	frames []*StackFrame
	size   int
}

func NewCallStack() CallStack {
	return CallStack{
		frames: make([]*StackFrame, initialFrameSize),
		size:   0,
	}
}

func (stack *CallStack) Push(closure *object.ClosureObject, bp, ip int) {
	stack.frames = ensureSpace(stack.frames, stack.size, 1)
	stack.frames[stack.size] = &StackFrame{closure, bp, ip}
	stack.size += 1
}

func (stack *CallStack) Peek() *StackFrame {
	ensureNotEmpty(stack.size)
	return stack.frames[stack.size-1]
}

func (stack *CallStack) Pop() *StackFrame {
	ensureNotEmpty(stack.size)
	stack.size -= 1
	frame := stack.frames[stack.size]
	stack.frames[stack.size] = nil
	return frame
}

func (stack *CallStack) Reset() {
	for i := stack.size - 1; i >= 0; i-- {
		stack.frames[i] = nil
	}
	stack.size = 0
}
