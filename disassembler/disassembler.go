package disassembler

import (
	"demo/bytecode"
	"demo/bytes"
	"fmt"
	"math"
)

func Disassemble(program []uint8) {

	index := 0

	constsSize := bytes.Get2Byte(program, index)
	index += 2

	fmt.Println("\nConstant Pool:")

	var constsCount = 0
	for i := uint16(0); i < constsSize; i++ {
		tag := bytecode.ConstTag(program[index])
		index++
		switch tag {
		case bytecode.TAG_INTEGER_CONST:
			value := int32(bytes.Get4Byte(program, index))
			fmt.Printf("    %04d - %d\n", constsCount, value)
			index += 4
		case bytecode.TAG_LONG_CONST:
			value := int64(bytes.Get8Byte(program, index))
			fmt.Printf("    %04d - %d\n", constsCount, value)
			index += 8
		case bytecode.TAG_FLOAT_CONST:
			value := math.Float32frombits(bytes.Get4Byte(program, index))
			fmt.Printf("    %04d - %f\n", constsCount, value)
			index += 4
		case bytecode.TAG_DOUBLE_CONST:
			value := math.Float64frombits(bytes.Get8Byte(program, index))
			fmt.Printf("    %04d - %f\n", constsCount, value)
			index += 8
		case bytecode.TAG_STRING_CONST:
			length := int(bytes.Get2Byte(program, index))
			index += 2
			bytes := make([]uint8, length)
			copy(bytes, program[index:index+length])
			fmt.Printf("    %04d - %s\n", constsCount, string(bytes))
			index += length
		case bytecode.TAG_FUN_CONST:
			ptr := bytes.Get2Byte(program, index)
			arity := bytes.Get1Byte(program, index+4)
			fmt.Printf("    %04d - Fun{ptr: %d, arity: %d}\n", constsCount, ptr, arity)
			index += 3
		default:
			panic(fmt.Sprintf("Unknown constant tag: %d", tag))
		}
		constsCount++
	}

	fmt.Printf("\nGlobals size: %d.\n", bytes.Get2Byte(program, index))
	index += 2

	fmt.Println("\nBytecode:")

	startIndex := index

	programSize := len(program)
	for index < programSize {

		opcode := program[index]
		instruction := instructions[opcode]
		index++

		fmt.Printf("    %04d - %s", (index-1)-startIndex, instruction.name)

		if opcode == uint8(bytecode.OP_CLOSURE) {

			freeCount := bytes.Get2Byte(program, index)
			fmt.Printf(" %d\n", freeCount)
			index += 2

			for i := 0; i < int(freeCount); i++ {
				isLocal := toBool(bytes.Get1Byte(program, index))
				freeIndex := bytes.Get2Byte(program, index+1)
				fmt.Printf("        [%d]: isLocal = %t\n", freeIndex, isLocal)
				index += 3
			}

			continue
		}

		max := len(instruction.operands)
		for i := 0; i < max; i++ {
			operandSize := instruction.operands[i]
			switch operandSize {
			case -2:
				fmt.Printf(" %d", int16(bytes.Get2Byte(program, index)))
			case -1:
				fmt.Printf(" %d", int8(bytes.Get1Byte(program, index)))
			case 1:
				fmt.Printf(" %d", bytes.Get1Byte(program, index))
			case 2:
				fmt.Printf(" %d", bytes.Get2Byte(program, index))
			default:
				panic("Unhandled operand size")
			}
			index += int(abs(operandSize))
		}

		fmt.Println()
	}

	fmt.Println()
}

func abs(value int8) int8 {
	if value < 0 {
		return -value
	}
	return value
}

func toBool(value uint8) bool {
	if value == 0 {
		return false
	}
	return true
}
