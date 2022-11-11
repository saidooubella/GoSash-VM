package vm

import (
	"demo/bytecode"
	"demo/bytes"
	"demo/disassembler"
	"demo/object"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"
)

const (
	FLAG_NONE uint8 = iota
	FLAG_TIMED
	FLAG_DEBUG
)

type VM struct {
	constants []object.Object
	bytecode  []uint8
	running   bool
	globals   []object.Object
	frames    CallStack
	stack     ObjectStack
	flag      uint8
}

func New(program []uint8, flag uint8) *VM {

	if flag == FLAG_DEBUG {
		disassembler.Disassemble(program)
	}

	constants, globals, program := createConstants(program)

	return &VM{
		constants: constants,
		bytecode:  program,
		stack:     NewObjectStack(),
		frames:    NewCallStack(),
		globals:   make([]object.Object, globals),
		flag:      flag,
	}
}

func createConstants(program []uint8) ([]object.Object, int, []uint8) {

	var index int = 0

	constsSize := bytes.Get2Byte(program, index)
	constants := make([]object.Object, constsSize)

	index += 2

	for i := uint16(0); i < constsSize; i++ {
		tag := bytecode.ConstTag(program[index])
		index++
		switch tag {
		case bytecode.TAG_INTEGER_CONST:
			value := int32(bytes.Get4Byte(program, index))
			constants[i] = object.NewInt(value)
			index += 4
		case bytecode.TAG_LONG_CONST:
			value := int64(bytes.Get8Byte(program, index))
			constants[i] = object.NewLong(value)
			index += 8
		case bytecode.TAG_FLOAT_CONST:
			value := bytes.Get4Byte(program, index)
			constants[i] = object.NewFloat(math.Float32frombits(value))
			index += 4
		case bytecode.TAG_DOUBLE_CONST:
			value := bytes.Get8Byte(program, index)
			constants[i] = object.NewDouble(math.Float64frombits(value))
			index += 8
		case bytecode.TAG_STRING_CONST:
			length := int(bytes.Get2Byte(program, index))
			index += 2
			bytes := make([]uint8, length)
			copy(bytes, program[index:index+length])
			constants[i] = object.NewUTF8(bytes)
			index += length
		case bytecode.TAG_FUN_CONST:
			ptr := bytes.Get2Byte(program, index)
			arity := bytes.Get1Byte(program, index+2)
			constants[i] = object.NewFunPtr(ptr, arity)
			index += 3
		default:
			panic(fmt.Sprintf("Unknown constant tag: %d", tag))
		}
	}

	globalsSize := bytes.Get2Byte(program, index)
	index += 2

	program = program[index:]
	bytecode := make([]uint8, len(program))
	copy(bytecode, program)

	return constants, int(globalsSize), bytecode
}

func (vm *VM) Start() error {

	vm.reset()

	start := time.Now()

	for vm.running {
		opcode := bytecode.OpCode(vm.fetchU8())
		if err := vm.execute(opcode); err != nil {
			return err
		}
	}

	end := time.Now()

	if vm.stack.Size() > 0 {
		return errors.New("Stack is not empty")
	}

	if vm.flag == FLAG_TIMED {
		fmt.Printf("Execution took: %d ms\n", end.UnixMilli()-start.UnixMilli())
	}

	return nil
}

func (vm *VM) reset() {
	vm.running = true
	vm.frames.Reset()
	vm.stack.Reset()
	vm.frames.Push(nil, 0, 0)
}

func (vm *VM) fetchU8() uint8 {
	vm.frames.Peek().ip += 1
	return vm.bytecode[vm.frames.Peek().ip-1]
}

func (vm *VM) fetchU16() uint16 {
	vm.frames.Peek().ip += 2
	return (uint16(vm.bytecode[vm.frames.Peek().ip-2])&0xFF)<<8 |
		(uint16(vm.bytecode[vm.frames.Peek().ip-1]) & 0xFF)
}

func (vm *VM) fetchI16() int16 {
	vm.frames.Peek().ip += 2
	return (int16(vm.bytecode[vm.frames.Peek().ip-2])&0xFF)<<8 |
		(int16(vm.bytecode[vm.frames.Peek().ip-1]) & 0xFF)
}

func (vm *VM) currentClosure() *object.ClosureObject {
	closure := vm.frames.Peek().closure
	if closure == nil {
		panic("current frame closure is nil where it shouldn't")
	}
	return closure
}

func (vm *VM) execute(opcode bytecode.OpCode) error {

	if vm.flag == FLAG_DEBUG {
		ip := vm.frames.Peek().ip - 1
		fmt.Printf("[ip: %d]: ", ip)
	}

	switch bytecode.OpCode(opcode) {
	case bytecode.OP_HALT:
		vm.running = false
	case bytecode.OP_PUSH_FALSE:
		vm.stack.Push(object.BoolFalse)
	case bytecode.OP_PUSH_TRUE:
		vm.stack.Push(object.BoolTrue)
	case bytecode.OP_POP:
		vm.stack.Pop()
	case bytecode.OP_CONSTANT:
		vm.stack.Push(vm.constants[vm.fetchU16()])
	case bytecode.OP_NOT:
		vm.stack.Push(vm.stack.Pop().(*object.BoolObject).Inv())
	case bytecode.OP_GOTO:
		vm.frames.Peek().ip += int(vm.fetchI16())
	case bytecode.OP_CALL:
		switch callable := vm.stack.Pop().(type) {
		case *object.FunPtrObject:
			vm.frames.Push(nil, vm.stack.Size()-int(callable.Arity), int(callable.Ptr))
		case *object.ClosureObject:
			vm.frames.Push(callable, vm.stack.Size()-int(callable.Func.Arity), int(callable.Func.Ptr))
		default:
			panic(fmt.Sprintf("Unknown callable type %T", callable))
		}
	case bytecode.OP_RETURN:
		frame := vm.frames.Pop()
		result := vm.stack.Pop()
		for frame.bp < vm.stack.Size() {
			vm.stack.Pop()
		}
		vm.stack.Push(result)
	case bytecode.OP_CONCAT:
		count := vm.fetchU16()
		var builder strings.Builder
		for index := uint16(0); index < count; index++ {
			builder.WriteString(vm.stack.Pop().(*object.StringObject).Value)
		}
		vm.stack.Push(object.NewString(builder.String()))
	case bytecode.OP_GOTO_DROP_FALSE:
		condition := vm.stack.Pop().(*object.BoolObject)
		target := int(vm.fetchI16())
		if !condition.Value {
			vm.frames.Peek().ip += target
		}
	case bytecode.OP_GOTO_DROP_TRUE:
		condition := vm.stack.Pop().(*object.BoolObject)
		target := int(vm.fetchI16())
		if condition.Value {
			vm.frames.Peek().ip += target
		}
	case bytecode.OP_GOTO_FALSE_OR_DROP:
		condition := vm.stack.Peek().(*object.BoolObject)
		target := int(vm.fetchI16())
		if condition.Value {
			vm.stack.Pop()
		} else {
			vm.frames.Peek().ip += target
		}
	case bytecode.OP_GOTO_TRUE_OR_DROP:
		condition := vm.stack.Peek().(*object.BoolObject)
		target := int(vm.fetchI16())
		if condition.Value {
			vm.frames.Peek().ip += target
		} else {
			vm.stack.Pop()
		}
	case bytecode.OP_SET_GLOBAL:
		index := vm.fetchU16()
		vm.globals[index] = vm.stack.Peek()
	case bytecode.OP_GET_GLOBAL:
		index := vm.fetchU16()
		vm.stack.Push(vm.globals[index])
	case bytecode.OP_GET_LOCAL:
		vm.stack.Push(vm.stack.Get(vm.frames.Peek().bp + int(vm.fetchU16())))
	case bytecode.OP_SET_LOCAL:
		vm.stack.Set(vm.frames.Peek().bp+int(vm.fetchU16()), vm.stack.Peek())
	case bytecode.OP_GET_FREE:
		currClosure := vm.currentClosure()
		vm.stack.Push(*currClosure.FreeValues[vm.fetchU16()].Pointer)
	case bytecode.OP_SET_FREE:
		currClosure := vm.currentClosure()
		currClosure.FreeValues[vm.fetchU16()].Pointer = vm.stack.PeekPtr()
	case bytecode.OP_CLOSE_FREE:
		panic("Not implemented")
	case bytecode.OP_CLOSURE:
		freeCount := vm.fetchU16()
		fun := vm.stack.Pop().(*object.FunPtrObject)
		closure := object.NewClosure(fun, freeCount)
		for i := 0; i < int(freeCount); i++ {
			switch vm.fetchU8() {
			case 0:
				closure.FreeValues[i] = vm.currentClosure().FreeValues[vm.fetchU16()]
			case 1:
				pointer := vm.stack.GetPtr(vm.frames.Peek().bp + int(vm.fetchU16()))
				closure.FreeValues[i] = &object.FreeValue{Pointer: pointer, Value: nil}
			default:
				panic("Unexpected isLocal state")
			}
		}
		vm.stack.Push(closure)
	case bytecode.OP_PRINT:
		fmt.Println(vm.stack.Pop().String())
	case bytecode.OP_ADD:
		right, left := vm.stack.Pop(), vm.stack.Pop()
		switch {
		case left.Type() == object.IntType && right.Type() == object.IntType:
			vm.stack.Push(left.(*object.IntObject).Plus(right.(*object.IntObject)))
		case left.Type() == object.FloatType && right.Type() == object.FloatType:
			vm.stack.Push(left.(*object.FloatObject).Plus(right.(*object.FloatObject)))
		case left.Type() == object.LongType && right.Type() == object.LongType:
			vm.stack.Push(left.(*object.LongObject).Plus(right.(*object.LongObject)))
		case left.Type() == object.DoubleType && right.Type() == object.DoubleType:
			vm.stack.Push(left.(*object.DoubleObject).Plus(right.(*object.DoubleObject)))
		default:
			panic("Unexpected '+' types combination")
		}
	case bytecode.OP_SUB:
		right, left := vm.stack.Pop(), vm.stack.Pop()
		switch {
		case left.Type() == object.IntType && right.Type() == object.IntType:
			vm.stack.Push(left.(*object.IntObject).Sub(right.(*object.IntObject)))
		case left.Type() == object.FloatType && right.Type() == object.FloatType:
			vm.stack.Push(left.(*object.FloatObject).Sub(right.(*object.FloatObject)))
		case left.Type() == object.LongType && right.Type() == object.LongType:
			vm.stack.Push(left.(*object.LongObject).Sub(right.(*object.LongObject)))
		case left.Type() == object.DoubleType && right.Type() == object.DoubleType:
			vm.stack.Push(left.(*object.DoubleObject).Sub(right.(*object.DoubleObject)))
		default:
			panic("Unexpected '-' types combination")
		}
	case bytecode.OP_MUL:
		right, left := vm.stack.Pop(), vm.stack.Pop()
		switch {
		case left.Type() == object.IntType && right.Type() == object.IntType:
			vm.stack.Push(left.(*object.IntObject).Mul(right.(*object.IntObject)))
		case left.Type() == object.FloatType && right.Type() == object.FloatType:
			vm.stack.Push(left.(*object.FloatObject).Mul(right.(*object.FloatObject)))
		case left.Type() == object.LongType && right.Type() == object.LongType:
			vm.stack.Push(left.(*object.LongObject).Mul(right.(*object.LongObject)))
		case left.Type() == object.DoubleType && right.Type() == object.DoubleType:
			vm.stack.Push(left.(*object.DoubleObject).Mul(right.(*object.DoubleObject)))
		default:
			panic("Unexpected '*' types combination")
		}
	case bytecode.OP_DIV:
		right, left := vm.stack.Pop(), vm.stack.Pop()
		switch {
		case left.Type() == object.IntType && right.Type() == object.IntType:
			vm.stack.Push(left.(*object.IntObject).Div(right.(*object.IntObject)))
		case left.Type() == object.FloatType && right.Type() == object.FloatType:
			vm.stack.Push(left.(*object.FloatObject).Div(right.(*object.FloatObject)))
		case left.Type() == object.LongType && right.Type() == object.LongType:
			vm.stack.Push(left.(*object.LongObject).Div(right.(*object.LongObject)))
		case left.Type() == object.DoubleType && right.Type() == object.DoubleType:
			vm.stack.Push(left.(*object.DoubleObject).Div(right.(*object.DoubleObject)))
		default:
			panic("Unexpected '/' types combination")
		}
	case bytecode.OP_EQUALS:
		right, left := vm.stack.Pop(), vm.stack.Pop()
		switch {
		case left.Type() == object.IntType && right.Type() == object.IntType:
			vm.stack.Push(left.(*object.IntObject).Equals(right.(*object.IntObject)))
		case left.Type() == object.FloatType && right.Type() == object.FloatType:
			vm.stack.Push(left.(*object.FloatObject).Equals(right.(*object.FloatObject)))
		case left.Type() == object.LongType && right.Type() == object.LongType:
			vm.stack.Push(left.(*object.LongObject).Equals(right.(*object.LongObject)))
		case left.Type() == object.DoubleType && right.Type() == object.DoubleType:
			vm.stack.Push(left.(*object.DoubleObject).Equals(right.(*object.DoubleObject)))
		case left.Type() == object.StringType && right.Type() == object.StringType:
			vm.stack.Push(left.(*object.StringObject).Equals(right.(*object.StringObject)))
		case left.Type() == object.BoolType && right.Type() == object.BoolType:
			vm.stack.Push(left.(*object.BoolObject).Equals(right.(*object.BoolObject)))
		default:
			panic("Unexpected '==' types combination")
		}
	case bytecode.OP_LESS:
		right, left := vm.stack.Pop(), vm.stack.Pop()
		switch {
		case left.Type() == object.IntType && right.Type() == object.IntType:
			vm.stack.Push(left.(*object.IntObject).Less(right.(*object.IntObject)))
		case left.Type() == object.FloatType && right.Type() == object.FloatType:
			vm.stack.Push(left.(*object.FloatObject).Less(right.(*object.FloatObject)))
		case left.Type() == object.LongType && right.Type() == object.LongType:
			vm.stack.Push(left.(*object.LongObject).Less(right.(*object.LongObject)))
		case left.Type() == object.DoubleType && right.Type() == object.DoubleType:
			vm.stack.Push(left.(*object.DoubleObject).Less(right.(*object.DoubleObject)))
		default:
			panic("Unexpected '<' types combination")
		}
	case bytecode.OP_GREATER:
		right, left := vm.stack.Pop(), vm.stack.Pop()
		switch {
		case left.Type() == object.IntType && right.Type() == object.IntType:
			vm.stack.Push(left.(*object.IntObject).Greater(right.(*object.IntObject)))
		case left.Type() == object.FloatType && right.Type() == object.FloatType:
			vm.stack.Push(left.(*object.FloatObject).Greater(right.(*object.FloatObject)))
		case left.Type() == object.LongType && right.Type() == object.LongType:
			vm.stack.Push(left.(*object.LongObject).Greater(right.(*object.LongObject)))
		case left.Type() == object.DoubleType && right.Type() == object.DoubleType:
			vm.stack.Push(left.(*object.DoubleObject).Greater(right.(*object.DoubleObject)))
		default:
			panic("Unexpected '>' types combination")
		}
	case bytecode.OP_NEG:
		switch operand := vm.stack.Pop().(type) {
		case *object.IntObject:
			vm.stack.Push(operand.Neg())
		case *object.FloatObject:
			vm.stack.Push(operand.Neg())
		case *object.LongObject:
			vm.stack.Push(operand.Neg())
		case *object.DoubleObject:
			vm.stack.Push(operand.Neg())
		default:
			panic("Unexpected negation type")
		}
	default:
		panic(fmt.Sprintf("Unknown opcode: %d", opcode))
	}

	if vm.flag == FLAG_DEBUG {
		vm.stack.Debug()
	}

	return nil
}
