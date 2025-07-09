package internal

import (
	"fmt"
)

func (i *ArithmeticInstruction) SimulateRegRM(simulator *Simulator) {
	prevFlags := simulator.flags
	var dest, src uint8 = 0, 0
	if i.isDestination() {
		dest = i.reg
		src = i.rm
	} else {
		dest = i.rm
		src = i.reg
	}
	previous := simulator.register[dest]
	value := simulator.register[dest]
	switch i.op {
	case SUB_REG_MEM:
		value -= simulator.register[src]
	case ADD_REG_MEM:
		value += simulator.register[src]
	case CMP_REG_MEM:
		value = simulator.register[dest] - simulator.register[src]
	}

	if value == 0 {
		simulator.flags.Set(FLAGS_ZERO, true)
	} else {
		simulator.flags.Set(FLAGS_ZERO, false)
	}

	if (uint16(value) >> 15) == 1 {
		simulator.flags.Set(FLAGS_SIGN, true)
	} else {
		simulator.flags.Set(FLAGS_SIGN, false)
	}

	switch i.op {
	case SUB_REG_MEM, ADD_REG_MEM:
		simulator.register[dest] = value
	case CMP_REG_MEM:
	}

	current := simulator.register[dest]
	curFlags := simulator.flags
	regName := REGISTERS_NAME[dest]
	operationStr, _ := i.Disassemble(&simulator.clocks)
	fmt.Printf("%s; %s:0x%x->0x%x (%d); flags:%v->%v; ", operationStr, regName, uint16(previous), uint16(current), current, prevFlags.String(), curFlags.String())
}
func (i *ArithmeticInstruction) SimulateImmediate(simulator *Simulator) {
	prevFlags := simulator.flags
	previous := simulator.register[i.rm]
	value := simulator.register[i.rm]

	switch i.op {
	case SUB_IMMEDIATE_RM:
		value -= i.data
	case ADD_IMMEDIATE_RM:
		value += i.data
	case CMP_IMMEDIATE_RM:
		value = value - i.data
	}

	if value == 0 {
		simulator.flags.Set(FLAGS_ZERO, true)
	} else {
		simulator.flags.Set(FLAGS_ZERO, false)
	}

	if (value << 15) == 1 {
		simulator.flags.Set(FLAGS_SIGN, true)
	} else {
		simulator.flags.Set(FLAGS_SIGN, false)
	}

	switch i.op {
	case SUB_IMMEDIATE_RM, ADD_IMMEDIATE_RM:
		simulator.register[i.rm] = value
	case CMP_IMMEDIATE_RM:
	}
	current := simulator.register[i.rm]
	regName := REGISTERS_NAME[i.rm]
	curFlags := simulator.flags
	operationStr, _ := i.Disassemble(&simulator.clocks)
	fmt.Printf("%s; %s:0x%x->0x%x (%d); flags:%v->%v; ", operationStr, regName, uint16(previous), uint16(current), current, prevFlags, curFlags)
}
