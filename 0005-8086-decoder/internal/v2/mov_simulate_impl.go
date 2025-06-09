package internal

import "fmt"

func (i *MovInstruction) getSimulateFuncMap() SimulateFuncTable {
	return SimulateFuncTable{
		// MOV_ACCULUMATOR_FROM_TO_MEMORY:   i.decodeMovAccumulatorFromToMemory,
		MOV_IMMEDIATE_TO_REGISTER: i.SimulateMovImmidiateToRegister,
		// MOV_IMMEDIATE_TO_REGISTER_MEMORY: i.decodeMovImmediateToRegisterMemory,
		MOV_REGISTER_FROM_TO_MEMORY: i.SimulateMovRMFromToRegister,
	}
}

func (i *MovInstruction) SimulateMovImmidiateToRegister(sim *Simulator) {
	previous := sim.memory[i.reg]
	sim.memory[i.reg] = i.data
	current := sim.memory[i.reg]
	regName := REGISTERS_NAME[i.reg]
	operationStr, _ := i.Disassemble()
	fmt.Printf("%s; %s:0x%x->0x%x (%v); ", operationStr, regName, uint16(previous), uint16(current), current)
}

func (i *MovInstruction) SimulateMovRMFromToRegister(sim *Simulator) {
	// only the non-memory mov are implemented
	var dest, src uint8 = 0, 0
	if i.isDestination() {
		dest = i.reg
		src = i.rm
	} else {
		dest = i.rm
		src = i.reg
	}
	previous := sim.memory[dest]
	sim.memory[dest] = sim.memory[src]
	current := sim.memory[dest]
	regName := REGISTERS_NAME[dest]
	operationStr, _ := i.Disassemble()
	fmt.Printf("%s; %s:0x%x->0x%x (%v); ", operationStr, regName, uint16(previous), uint16(current), current)
}
