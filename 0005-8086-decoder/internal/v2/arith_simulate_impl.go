package internal

import (
	"fmt"
)

func (i *ArithmeticInstruction) SimulateRegRM(memory *Memory, flags *Flags) {
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

	if (value << 15) == 1 {
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
	regName := REGISTERS_NAME[dest]
	operationStr, _ := i.Disassemble()
	fmt.Printf("%s ; %s:0x%x->0x%x\n", operationStr, regName, previous, current)
}
func (i *ArithmeticInstruction) SimulateImmediate(memory *Memory, flags *Flags) {
	previous := memory[i.reg]
	value := memory[i.reg]

	switch i.op {
	case SUB_REG_MEM:
		value -= i.data
	case ADD_REG_MEM:
		value += i.data
	case CMP_REG_MEM:
		value = memory[i.reg] - i.data
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
	case SUB_REG_MEM, ADD_REG_MEM:
		memory[i.reg] = value
	case CMP_REG_MEM:
	}

	current := memory[i.reg]
	regName := REGISTERS_NAME[i.reg]
	operationStr, _ := i.Disassemble()
	fmt.Printf("%s ; %s:0x%x->0x%x\n", operationStr, regName, previous, current)
}
