package internal

import (
	"fmt"
	"log"
)

func (i *MovInstruction) getSimulateFuncMap() SimulateFuncTable {
	return SimulateFuncTable{
		// MOV_ACCULUMATOR_FROM_TO_MEMORY:   i.decodeMovAccumulatorFromToMemory,
		MOV_IMMEDIATE_TO_REGISTER:        i.SimulateMovImmidiateToRegister,
		MOV_IMMEDIATE_TO_REGISTER_MEMORY: i.SimulateMovImmideateToRegisterMemory,
		MOV_REGISTER_FROM_TO_MEMORY:      i.SimulateMovRMFromToRegister,
	}
}

func (i *MovInstruction) SimulateMovImmideateToRegisterMemory(sim *Simulator) {
	if i.isEffectiveAddress() {
		memoryAddress := i.handleEffectiveAddressCalculation(sim)
		sim.memory[memoryAddress] = uint8(i.data)
		sim.memory[memoryAddress+1] = uint8(i.data >> 8)
		operationStr, _ := i.Disassemble(&sim.clocks)
		fmt.Printf("%s; ", operationStr)
	} else {
		previous := sim.register[i.rm]
		sim.register[i.rm] = i.data
		current := sim.register[i.rm]
		regName := REGISTERS_NAME[i.rm]
		operationStr, _ := i.Disassemble(&sim.clocks)
		fmt.Printf("%s; %s:0x%x->0x%x (%v); ", operationStr, regName, uint16(previous), uint16(current), current)
	}
}
func (i *MovInstruction) SimulateMovImmidiateToRegister(sim *Simulator) {
	previous := sim.register[i.reg]
	sim.register[i.reg] = i.data
	current := sim.register[i.reg]
	regName := REGISTERS_NAME[i.reg]
	operationStr, _ := i.Disassemble(&sim.clocks)
	fmt.Printf("%s; %s:0x%x->0x%x (%v); ", operationStr, regName, uint16(previous), uint16(current), current)
}

func (i *MovInstruction) SimulateMovRMFromToRegister(sim *Simulator) {
	// only the non-memory mov are implemented
	if i.isEffectiveAddress() {
		memoryAddress := i.handleEffectiveAddressCalculation(sim)
		if i.isDestination() {
			// little endian (least significant byte is store at the lowest address)
			sim.register[i.reg] = int16(uint16(sim.memory[memoryAddress]) | uint16(sim.memory[memoryAddress+1])<<8)
		} else {
			sim.memory[memoryAddress] = uint8(sim.register[i.reg])
			sim.memory[memoryAddress+1] = uint8(sim.register[i.reg] >> 8)
		}
		operationStr, _ := i.Disassemble(&sim.clocks)
		fmt.Printf("%s; ", operationStr)
	} else {
		var dest, src uint8 = 0, 0
		if i.isDestination() {
			dest = i.reg
			src = i.rm
		} else {
			dest = i.rm
			src = i.reg
		}
		previous := sim.register[dest]
		sim.register[dest] = sim.register[src]
		current := sim.register[dest]
		regName := REGISTERS_NAME[dest]
		operationStr, _ := i.Disassemble(&sim.clocks)
		fmt.Printf("%s; %s:0x%x->0x%x (%v); ", operationStr, regName, uint16(previous), uint16(current), current)
	}
}

func (i *MovInstruction) getEffectiveAddressCalcFunc() map[byte]func(sim *Simulator) int16 {
	return map[byte]func(sim *Simulator) int16{
		0b000: func(sim *Simulator) int16 {
			return sim.register[0b011] + sim.register[0b110] // bx + si
		},
		0b001: func(sim *Simulator) int16 {
			return sim.register[0b011] + sim.register[0b111] // bx + di
		},
		0b010: func(sim *Simulator) int16 {
			return sim.register[0b101] + sim.register[0b110] // bp + si
		},
		0b011: func(sim *Simulator) int16 {
			return sim.register[0b101] + sim.register[0b111] // bp + di
		},
		0b100: func(sim *Simulator) int16 {
			return sim.register[0b110]
		},
		0b101: func(sim *Simulator) int16 {
			return sim.register[0b111]
		},
		0b110: func(sim *Simulator) int16 {
			return sim.register[0b101]
		},
		0b111: func(sim *Simulator) int16 {
			return sim.register[0b011]
		},
	}
}
func (i *MovInstruction) handleEffectiveAddressCalculation(sim *Simulator) int16 {
	if i.isDirectAccess() {
		return i.handleDisplacepment()
	}

	eac, found := i.getEffectiveAddressCalcFunc()[i.rm]
	if !found {
		log.Printf("effective address calculation function for value %b not found", i.rm)
		return 0
	}
	value := eac(sim)
	if i.isDisplacement() {
		value += i.handleDisplacepment()
	}
	return value
}

// func (i *MovInstruction) handleWriteToMemory(sim *Simulator,)
