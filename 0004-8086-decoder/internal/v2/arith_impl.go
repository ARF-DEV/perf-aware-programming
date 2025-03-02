package internal

import "fmt"

type ArithmeticInstruction struct {
	d, w, mod, reg, rm, s uint8
	lo, hi                uint8
	data                  int16
	op                    OpMode
}

func (i *ArithmeticInstruction) String() string {
	return fmt.Sprintf("op:%d\nd:%d\tw:%d\tmod:%02b\treg:%03b\trm:%03b\nlo:%08b\thi:%08b\ndata:%d", i.op, i.d, i.w, i.mod, i.reg, i.rm, i.lo, i.hi, i.data)
}
func (i *ArithmeticInstruction) Disassemble() (string, error) {
	decode, found := i.getDecoderFuncMap()[i.op]
	if !found {
		return "", fmt.Errorf("error: operation not implemented")
	}
	return decode(), nil
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

func (i *ArithmeticInstruction) decodeRM() string {
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

	decode = fmt.Sprintf("%s %s, %s", decode, dst, src)
	return decode
}

func (i *ArithmeticInstruction) decodeImmediate() string {
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
	decode = i.getOperationDecodeMap()[i.op]
	dst = rmStr
	src = fmt.Sprintf("%d", i.data)
	// if i.isWord() {
	// 	src = "word " + src
	// } else {
	// 	src = "byte " + src
	// }

	decode = fmt.Sprintf("%s %s, %s", decode, dst, src)
	return decode
}
func (i *ArithmeticInstruction) decodeAccumulator() string {
	regStr := RegisterTab.Get(0b11, i.w, i.reg)

	var dst, src string
	if !i.isDestination() {
		dst = regStr
		src = fmt.Sprintf("%d", i.handleDisplacepment())
	} else {
		src = regStr
		dst = fmt.Sprintf("%d", i.handleDisplacepment())
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
	// fmt.Println(disp)
	return disp
}
