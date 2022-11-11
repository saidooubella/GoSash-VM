// AUTO GENERATED - DO NOT MODIFY
package disassembler

type instruction struct {
	name string
	operands []int8
}

var instructions = [...]instruction{
	{"op_goto_false_or_drop", []int8{-2}},
	{"op_goto_true_or_drop", []int8{-2}},
	{"op_goto_drop_false", []int8{-2}},
	{"op_goto_drop_true", []int8{-2}},
	{"op_goto", []int8{-2}},
	{"op_close_free", []int8{}},
	{"op_push_false", []int8{}},
	{"op_push_true", []int8{}},
	{"op_add", []int8{}},
	{"op_sub", []int8{}},
	{"op_mul", []int8{}},
	{"op_div", []int8{}},
	{"op_neg", []int8{}},
	{"op_get_global", []int8{2}},
	{"op_set_global", []int8{2}},
	{"op_get_free", []int8{2}},
	{"op_set_free", []int8{2}},
	{"op_get_local", []int8{2}},
	{"op_set_local", []int8{2}},
	{"op_return", []int8{}},
	{"op_call", []int8{}},
	{"op_greater", []int8{}},
	{"op_equals", []int8{}},
	{"op_less", []int8{}},
	{"op_not", []int8{}},
	{"op_concat", []int8{2}},
	{"op_print", []int8{}},
	{"op_constant", []int8{2}},
	{"op_closure", []int8{2}},
	{"op_pop", []int8{}},
	{"op_halt", []int8{}},
}
