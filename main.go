package main

import (
	"demo/vm"
	"fmt"
	"os"
)

func main() {

	name := "C:\\Users\\Said\\Desktop\\main.shp"
	program, err := os.ReadFile(name)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// disassembler.Disassemble(program)

	vm := vm.New(program, vm.FLAG_DEBUG)

	if err := vm.Start(); err != nil {
		fmt.Printf("Error: %s.\n", err.Error())
	}
}
