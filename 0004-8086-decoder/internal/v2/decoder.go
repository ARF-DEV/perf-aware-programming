package internal

import (
	"bufio"
	"os"
	"strings"
)

type InstructionDecoder struct {
	buf     []byte
	curIdx  int64
	nextIdx int64
	builder strings.Builder

	instructionFuncs  map[byte]InstructionFunc
	instructionFuncss map[int]map[byte]InstructionFunc
}

func NewDecoder(filename string) *InstructionDecoder {
	ins := InstructionDecoder{}
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	info, err := f.Stat()
	if err != nil {
		panic(err)
	}
	r := bufio.NewReader(f)
	ins.buf = make([]byte, info.Size())
	_, err = r.Read(ins.buf)
	if err != nil {
		panic(err)
	}

	ins.builder = strings.Builder{}
	ins.builder.WriteString("bits 16\n\n")
	return &ins
}

func (i *InstructionDecoder) MovInstruction() InstructionStatement {
	stmt := MovInstruction{op: MOV_REGISTER_FROM_TO_MEMORY}
	stmt.w = i.getBits(0, 1)
	stmt.d = i.getBits(1, 1)

	i.Next()
	stmt.mod = i.getBits(6, 2)
	stmt.reg = i.getBits(3, 3)
	stmt.rm = i.getBits(0, 3)
	switch stmt.mod {
	case 0b00:
		if stmt.rm == 0b110 {
			i.Next()
			stmt.lo = int8(i.CurrentByte())
			i.Next()
			stmt.hi = int8(i.CurrentByte())
		}
	case 0b01:
		i.Next()
		stmt.lo = int8(i.CurrentByte())
	case 0b10:
		i.Next()
		stmt.lo = int8(i.CurrentByte())
		i.Next()
		stmt.hi = int8(i.CurrentByte())
	default:
	}

	return &stmt
}

func (i *InstructionDecoder) MovAccumulator() InstructionStatement {
	stmt := MovInstruction{op: MOV_ACCULUMATOR_FROM_TO_MEMORY}
	stmt.w = i.getBits(0, 1)
	stmt.d = i.getBits(1, 1)

	i.Next()
	stmt.reg = 0b000
	switch stmt.w {
	case 0b0:
		i.Next()
		stmt.lo = int8(i.CurrentByte())
	case 0b1:
		i.Next()
		stmt.lo = int8(i.CurrentByte())
		i.Next()
		stmt.hi = int8(i.CurrentByte())
	default:
	}

	return &stmt
}
func (i *InstructionDecoder) MovIRegInstruction() InstructionStatement {
	stmt := MovInstruction{op: MOV_IMMEDIATE_TO_REGISTER}
	stmt.reg = i.getBits(0, 3)
	stmt.w = i.getBits(3, 1)

	// regStr := RegisterMap[reg|w<<3|1<<4]
	switch stmt.w {
	case 0b1:
		var b []uint16
		i.Next()
		b = append(b, uint16(i.CurrentByte()))
		i.Next()
		b = append(b, uint16(i.CurrentByte()))

		stmt.data = int16((b[1] << 8) | b[0])
	case 0b0:
		i.Next()
		// convert properly to int8 first then convert it to 16 bit
		stmt.data = int16(int8(i.CurrentByte()))
	}
	return &stmt
}
func (i *InstructionDecoder) MovIRMInstruction() InstructionStatement {
	// fmt.Println("MOV RRM")
	stmt := MovInstruction{op: MOV_IMMEDIATE_TO_REGISTER_MEMORY}
	w := i.getBits(0, 1)

	i.Next()
	stmt.mod = i.getBits(6, 2)

	stmt.rm = i.getBits(0, 3)
	// rmStr := RegisterMap[rm|(w<<3)|byte(RtR)<<4]

	switch stmt.mod {
	case 0b00:
		if stmt.rm == 0b110 {
			i.Next()
			stmt.lo = int8(i.CurrentByte())
			i.Next()
			stmt.hi = int8(i.CurrentByte())
		}
	case 0b01:
		i.Next()
		stmt.lo = int8(i.CurrentByte())
	case 0b10:
		i.Next()
		stmt.lo = int8(i.CurrentByte())
		i.Next()
		stmt.hi = int8(i.CurrentByte())
	}

	switch w {
	case 0b1:
		var b []uint16
		i.Next()
		b = append(b, uint16(i.CurrentByte()))
		i.Next()
		b = append(b, uint16(i.CurrentByte()))
		stmt.data = int16(b[1]<<8 | b[0])
	case 0b0:
		i.Next()
		stmt.data = int16(i.CurrentByte())
	}
	return &stmt
}
func (i *InstructionDecoder) Next() bool {
	i.curIdx = i.nextIdx
	i.nextIdx++
	return int(i.curIdx) < len(i.buf)
}
func (i *InstructionDecoder) CurrentByte() byte {
	return i.buf[i.curIdx]
}
func (i *InstructionDecoder) NextByte() byte {
	return i.buf[i.nextIdx]
}

func (i *InstructionDecoder) getBits(shift, nBits uint8) byte {
	return (i.CurrentByte() >> shift) & ((nBits << 1) - 1)
}
