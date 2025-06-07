package internal

import (
	"fmt"
)

func (i *ArithmeticInstruction) SimulateRegRM(memory *Memory, flags *Flags) {
	prevFlags := *flags
	var dest, src uint8 = 0, 0
	if i.isDestination() {
		dest = i.reg
		src = i.rm
	} else {
		dest = i.rm
		src = i.reg
	}
	previous := memory[dest]
	value := memory[dest]
	switch i.op {
	case SUB_REG_MEM:
		value -= memory[src]
	case ADD_REG_MEM:
		value += memory[src]
	case CMP_REG_MEM:
		value = memory[dest] - memory[src]
	}

	if value == 0 {
		flags.Set(FLAGS_ZERO, true)
	} else {
		flags.Set(FLAGS_ZERO, false)
	}

	if (uint16(value) >> 15) == 1 {
		flags.Set(FLAGS_SIGN, true)
	} else {
		flags.Set(FLAGS_SIGN, false)
	}

	switch i.op {
	case SUB_REG_MEM, ADD_REG_MEM:
		memory[dest] = value
	case CMP_REG_MEM:
	}

	current := memory[dest]
	curFlags := *flags
	regName := REGISTERS_NAME[dest]
	operationStr, _ := i.Disassemble()
	fmt.Printf("%s ; %s:0x%x->0x%x (%d); flags:%v->%v\n", operationStr, regName, uint16(previous), uint16(current), current, prevFlags.String(), curFlags.String())
}
func (i *ArithmeticInstruction) SimulateImmediate(memory *Memory, flags *Flags) {
	prevFlags := *flags
	previous := memory[i.rm]
	value := memory[i.rm]

	switch i.op {
	case SUB_IMMEDIATE_RM:
		value -= i.data
	case ADD_IMMEDIATE_RM:
		value += i.data
	case CMP_IMMEDIATE_RM:
		value = value - i.data
	}

	if value == 0 {
		flags.Set(FLAGS_ZERO, true)
	} else {
		flags.Set(FLAGS_ZERO, false)
	}

	if (value << 15) == 1 {
		flags.Set(FLAGS_SIGN, true)
	} else {
		flags.Set(FLAGS_SIGN, false)
	}

	switch i.op {
	case SUB_IMMEDIATE_RM, ADD_IMMEDIATE_RM:
		memory[i.rm] = value
	case CMP_IMMEDIATE_RM:
	}
	current := memory[i.rm]
	regName := REGISTERS_NAME[i.rm]
	curFlags := *flags
	operationStr, _ := i.Disassemble()
	fmt.Printf("%s ; %s:0x%x->0x%x (%d); flags:%v->%v\n", operationStr, regName, uint16(previous), uint16(current), current, prevFlags, curFlags)
}
