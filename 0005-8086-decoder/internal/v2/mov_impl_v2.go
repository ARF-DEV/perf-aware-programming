package internal

import "fmt"

type MovInstruction struct {
	d, w, mod, reg, rm uint8
	lo, hi             uint8
	data               int16
	op                 OpMode
}

func (i *MovInstruction) String() string {
	return fmt.Sprintf("op:%d\nd:%d\tw:%d\tmod:%02b\treg:%03b\trm:%03b\nlo:%08b\thi:%08b\ndata:%d", i.op, i.d, i.w, i.mod, i.reg, i.rm, i.lo, i.hi, i.data)
}
func (i *MovInstruction) Disassemble() (string, error) {
	decode, found := i.getDecoderFuncMap()[i.op]
	if !found {
		return "", fmt.Errorf("error: operation not implemented")
	}
	return decode(), nil
}

func (i *MovInstruction) Simulate(mem *Memory) {
	// simple immediate to memory
	// TODO: do all mov operations
	previous := mem[i.reg]
	mem[i.reg] = i.data
	current := mem[i.reg]
	regName := REGISTERS_NAME[i.reg]
	operationStr, _ := i.Disassemble()
	fmt.Printf("%s ; %s:0x%x->0x%x\n", operationStr, regName, previous, current)
}

func (i *MovInstruction) isInstruction() {}

func (i *MovInstruction) getDecoderFuncMap() DecoderFuncTable {
	return DecoderFuncTable{
		MOV_ACCULUMATOR_FROM_TO_MEMORY:   i.decodeMovAccumulatorFromToMemory,
		MOV_IMMEDIATE_TO_REGISTER:        i.decodeMovImmediateToRegister,
		MOV_IMMEDIATE_TO_REGISTER_MEMORY: i.decodeMovImmediateToRegisterMemory,
		MOV_REGISTER_FROM_TO_MEMORY:      i.decodeMovRegisterFromToMemory,
	}
}

func (i *MovInstruction) decodeMovRegisterFromToMemory() string {
	regStr := RegisterTab.Get(0b11, i.w, i.reg)
	rmStr := RegisterTab.Get(i.mod, i.w, i.rm)
	if i.isEffectiveAddress() {
		if i.isDisplacement() {
			displacement := i.handleDisplacepment()
			if displacement != 0 {
				op := "+"
				if displacement < 0 {
					op = "-"
					displacement *= -1
				}
				rmStr = fmt.Sprintf("%s %s %d", rmStr, op, displacement)
			}
		} else {
			// handle if direct access
			if i.isDirectAccess() {
				rmStr = fmt.Sprintf("%d", int16(i.hi)<<8|int16(i.lo))
			}
		}
		rmStr = fmt.Sprintf("[%s]", rmStr)
	}

	var dst, src, decode string
	decode = "mov"
	if i.isDestination() {
		dst = regStr
		src = rmStr
	} else {
		src = regStr
		dst = rmStr
	}

	decode = fmt.Sprintf("%s %s, %s", decode, dst, src)
	return decode
}

func (i *MovInstruction) decodeMovImmediateToRegister() string {
	regStr := RegisterTab.Get(0b11, i.w, i.reg)
	return fmt.Sprintf("mov %s, %d", regStr, i.data)
}

func (i *MovInstruction) decodeMovImmediateToRegisterMemory() string {
	rmStr := RegisterTab.Get(i.mod, i.w, i.rm)

	if i.isEffectiveAddress() {
		if i.isDisplacement() {
			displacement := i.handleDisplacepment()
			if displacement != 0 {
				op := "+"
				if displacement < 0 {
					op = "-"
					displacement *= -1
				}
				rmStr = fmt.Sprintf("%s %s %d", rmStr, op, displacement)
			}
		} else {
			// handle if direct access
			if i.isDirectAccess() {
				rmStr = fmt.Sprintf("%d", i.handleDisplacepment())
			}
		}
		rmStr = fmt.Sprintf("[%s]", rmStr)
	}
	var dst, src, decode string
	decode = "mov"
	dst = rmStr
	src = fmt.Sprintf("%d", i.data)
	if i.isWord() {
		src = "word " + src
	} else {
		src = "byte " + src
	}

	decode = fmt.Sprintf("%s %s, %s", decode, dst, src)
	return decode
}
func (i *MovInstruction) decodeMovAccumulatorFromToMemory() string {
	regStr := RegisterTab.Get(0b11, i.w, i.reg)

	var dst, src string
	if !i.isDestination() {
		dst = regStr
		src = fmt.Sprintf("[%d]", i.handleDataInDisp())
	} else {
		src = regStr
		dst = fmt.Sprintf("[%d]", i.handleDataInDisp())
	}
	return fmt.Sprintf("mov %s, %s", dst, src)
}

func (i *MovInstruction) isDisplacement() bool {
	return !(i.mod == 0b11 || i.mod == 0b00)
}
func (i *MovInstruction) isDestination() bool {
	return i.d == 1
}
func (i *MovInstruction) isWord() bool {
	return i.w == 1
}
func (i *MovInstruction) isDirectAccess() bool {
	return i.rm == 0b110 && i.mod == 0b00
}
func (i *MovInstruction) isEffectiveAddress() bool {
	return i.mod != 0b11
}
func (i *MovInstruction) handleDisplacepment() int16 {
	var disp int16
	switch i.mod {
	case 0b01:
		disp = int16(int8(i.lo))
	case 0b10, 0b00:
		disp = int16(uint16(i.hi)<<8 | uint16(i.lo))
	}
	return disp
}

func (i *MovInstruction) handleDataInDisp() int16 {
	var srcInt int16
	switch i.w {
	case 0:
		srcInt = int16(int8(i.lo))
	case 1:
		srcInt = int16(uint16(i.hi)<<8 | uint16(i.lo))
	}
	return srcInt
}
