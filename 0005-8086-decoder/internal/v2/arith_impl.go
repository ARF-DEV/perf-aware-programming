package internal

import (
	"fmt"
	"log"
	"strings"
)

type ArithmeticInstruction struct {
	d, w, mod, reg, rm, s uint8
	lo, hi                uint8
	data                  int16
	op                    OpMode
}

func (i *ArithmeticInstruction) String() string {
	return fmt.Sprintf("op:%d\nd:%d\tw:%d\tmod:%02b\treg:%03b\trm:%03b\nlo:%08b\thi:%08b\ndata:%d", i.op, i.d, i.w, i.mod, i.reg, i.rm, i.lo, i.hi, i.data)
}
func (i *ArithmeticInstruction) Disassemble(clocks *int) (string, error) {
	decode, found := i.getDecoderFuncMap()[i.op]
	if !found {
		return "", fmt.Errorf("error: operation not implemented")
	}
	return decode(clocks), nil
}

func (i *ArithmeticInstruction) Simulate(simulator *Simulator) {
	sim, found := i.getSimulateFuncMap()[i.op]
	if !found {
		log.Printf("error: simulate function for op %d are not implemented", i.op)
		return
	}
	sim(simulator)
}

func (i *ArithmeticInstruction) isInstruction() {}

func (i *ArithmeticInstruction) getDecoderFuncMap() DecoderFuncTable {
	return DecoderFuncTable{
		ADD_REG_MEM:      i.decodeRM,
		ADD_IMMEDIATE_RM: i.decodeImmediate,
		ADD_ACC:          i.decodeAccumulator,
		SUB_REG_MEM:      i.decodeRM,
		SUB_IMMEDIATE_RM: i.decodeImmediate,
		SUB_ACC:          i.decodeAccumulator,
		CMP_REG_MEM:      i.decodeRM,
		CMP_IMMEDIATE_RM: i.decodeImmediate,
		CMP_ACC:          i.decodeAccumulator,
	}
}

func (i *ArithmeticInstruction) getSimulateFuncMap() SimulateFuncTable {
	return SimulateFuncTable{
		ADD_REG_MEM:      i.SimulateRegRM,
		ADD_IMMEDIATE_RM: i.SimulateImmediate,
		// ADD_ACC:          i.Simulate,
		SUB_REG_MEM:      i.SimulateRegRM,
		SUB_IMMEDIATE_RM: i.SimulateImmediate,
		// SUB_ACC:          i.decodeAccumulator,
		CMP_REG_MEM:      i.SimulateRegRM,
		CMP_IMMEDIATE_RM: i.SimulateImmediate,
		// CMP_ACC:          i.decodeAccumulator,
	}
}
func (i *ArithmeticInstruction) decodeRM(clocks *int) string {
	regStr := RegisterTab.Get(0b11, i.w, i.reg)
	rmStr := RegisterTab.Get(i.mod, i.w, i.rm)
	addClocks := i.estimateClocks()
	*clocks += addClocks
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
			if i.isDirectAccess() {
				rmStr = fmt.Sprintf("%d", int16(i.hi)<<8|int16(i.lo))
			}
		}
		rmStr = fmt.Sprintf("[%s]", rmStr)
	}

	var dst, src, decode string
	decode = i.getOperationDecodeMap()[i.op]
	if i.isDestination() {
		dst = regStr
		src = rmStr
	} else {
		src = regStr
		dst = rmStr
	}

	decode = fmt.Sprintf("%s %s, %s; Clocks += %d -> %d", decode, dst, src, addClocks, *clocks)
	return decode
}

func (i *ArithmeticInstruction) estimateClocks() int {
	clocks := 0
	rmStr := RegisterTab.Get(i.mod, i.w, i.rm)
	if i.isEffectiveAddress() {
		// displacement := i.handleDisplacepment()
		if strings.Contains(rmStr, "+") {
			switch rmStr {
			case "bp + di", "bx + si":
				clocks += 7
			case "bp + si", "bx + di":
				clocks += 8
			}
		} else if !i.isDirectAccess() {
			clocks += 5
		}
		if i.isDisplacement() && i.handleDisplacepment() != 0 {
			clocks += 4
		} else if i.isDirectAccess() {
			clocks += 6
		}
		if i.isDestination() {
			clocks += 9
		} else {
			clocks += 16
		}
	} else {
		clocks += 3
	}
	return clocks
}
func (i *ArithmeticInstruction) decodeImmediate(clocks *int) string {
	rmStr := RegisterTab.Get(i.mod, i.w, i.rm)
	// TBD
	*clocks += 4
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
		if i.isWord() {
			rmStr = "word " + rmStr
		} else {
			rmStr = "byte " + rmStr
		}
	}
	var dst, src, decode string
	decode = i.getOperationDecodeMap()[i.op]
	dst = rmStr
	src = fmt.Sprintf("%d", i.data)
	// if i.isEffectiveAddress() {
	// }

	decode = fmt.Sprintf("%s %s, %s; Clocks += 4 -> %d", decode, dst, src, *clocks)
	return decode
}
func (i *ArithmeticInstruction) decodeAccumulator(clocks *int) string {
	regStr := RegisterTab.Get(0b11, i.w, i.reg)

	var dst, src string
	if !i.isDestination() {
		dst = regStr
		src = fmt.Sprintf("%d", i.handleDataInDisp())
	} else {
		src = regStr
		dst = fmt.Sprintf("%d", i.handleDataInDisp())
	}
	return fmt.Sprintf("%s %s, %s", i.getOperationDecodeMap()[i.op], dst, src)
}
func (i *ArithmeticInstruction) getOperationDecodeMap() map[OpMode]string {
	return map[OpMode]string{
		ADD_REG_MEM:      "add",
		ADD_IMMEDIATE_RM: "add",
		ADD_ACC:          "add",
		SUB_REG_MEM:      "sub",
		SUB_IMMEDIATE_RM: "sub",
		SUB_ACC:          "sub",
		CMP_REG_MEM:      "cmp",
		CMP_IMMEDIATE_RM: "cmp",
		CMP_ACC:          "cmp",
	}
}

func (i *ArithmeticInstruction) isDisplacement() bool {
	return !(i.mod == 0b11 || i.mod == 0b00)
}
func (i *ArithmeticInstruction) isDestination() bool {
	return i.d == 1
}
func (i *ArithmeticInstruction) isWord() bool {
	return i.w == 1
}
func (i *ArithmeticInstruction) isDirectAccess() bool {
	return i.rm == 0b110 && i.mod == 0b00
}
func (i *ArithmeticInstruction) isEffectiveAddress() bool {
	return i.mod != 0b11
}
func (i *ArithmeticInstruction) handleDisplacepment() int16 {
	var disp int16
	switch i.mod {
	case 0b01:
		disp = int16(int8(i.lo))
	case 0b10, 0b00:
		disp = int16(uint16(i.hi)<<8 | uint16(i.lo))
	}
	return disp
}

func (i *ArithmeticInstruction) handleDataInDisp() int16 {
	var srcInt int16
	switch i.w {
	case 0:
		srcInt = int16(int8(i.lo))
	case 1:
		srcInt = int16(uint16(i.hi)<<8 | uint16(i.lo))
	}
	return srcInt
}
