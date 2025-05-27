package internal

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

type InstructionDecoder struct {
	buf        []byte
	curIdx     int64
	nextIdx    int64
	builder    strings.Builder
	statements []InstructionStatement

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
	ins.initMap()
	ins.builder = strings.Builder{}
	ins.builder.WriteString("bits 16\n\n")
	return &ins
}

func (i *InstructionDecoder) Decode() {
	// fmt.Printf("%08b\n", i.buf)
	for i.Next() {
		b := i.CurrentByte()
		// i.printCurrentByte()
		for j := 0; j < 8; j++ {
			ins, ok := i.instructionFuncss[8-j][b>>j]
			if !ok {
				continue
			}
			stmt := ins()
			// fmt.Println(stmt)
			// fmt.Println()
			i.statements = append(i.statements, stmt)
			break
		}
	}
}
func (i *InstructionDecoder) Disassemble(writer io.StringWriter) error {
	for _, stmt := range i.statements {
		v, err := stmt.Disassemble()
		if err != nil {
			return err
		}
		writer.WriteString(fmt.Sprintln(v))
		// fmt.Println(stmt)
		// fmt.Println()
	}
	return nil
}

func (i *InstructionDecoder) initMap() {
	i.instructionFuncs = map[byte]InstructionFunc{
		0b100010:  i.MovInstruction,
		0b1011:    i.MovIRegInstruction,
		0b1100011: i.MovIRMInstruction,
		// 0b00:      i.ArithRMAcc,
		// 0b100:     i.ArithImmediateInstruction,
		// 0b1010000: i.MovMToAcc,
		// 0b1010001: i.MovAccToM,

		// 0b000000:  i.IncRegToMem,
		// 0b001010:  i.IncRegToMem,
		// 0b001110:  i.IncRegToMem,
		// 0b100000:  i.ImmediateToRM,
		// 0b0000010: i.ImmediateToAcc,
		// 0b0010110: i.ImmediateToAcc,
		// 0b0011110: i.ImmediateToAcc,
	}
	i.instructionFuncss = map[int]map[byte]InstructionFunc{
		6: {
			0b100010: i.MovInstruction,
			0b101000: i.MovAccumulator,
			0b110001: i.MovIRMInstruction,
			// 0b000000: i.IncRegToMem,
			// 0b001010: i.IncRegToMem,
			// 0b001110: i.IncRegToMem,
			// 0b100000: i.ImmediateToRM,
		},
		4: {
			0b1011: i.MovIRegInstruction,
			0b0111: i.JumpLoop,
			0b1110: i.JumpLoop,
		},
		7: {
			// 0b1010001: i.MovAccumulator,
			// 0b1010000: i.MovAccumulator,
			// 0b0000010: i.ImmediateToAcc,
			// 0b0010110: i.ImmediateToAcc,
			// 0b0011110: i.ImmediateToAcc,
		},
		2: {
			0b00: i.ArithRMAcc,
		},
		3: {
			0b100: i.ArithImmediateInstruction,
		},
	}
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
			stmt.lo = (i.CurrentByte())
			i.Next()
			stmt.hi = (i.CurrentByte())
		}
	case 0b01:
		i.Next()
		stmt.lo = (i.CurrentByte())
	case 0b10:
		i.Next()
		stmt.lo = (i.CurrentByte())
		i.Next()
		stmt.hi = (i.CurrentByte())
	default:
	}

	return &stmt
}

func (i *InstructionDecoder) MovAccumulator() InstructionStatement {
	stmt := MovInstruction{op: MOV_ACCULUMATOR_FROM_TO_MEMORY}
	stmt.w = i.getBits(0, 1)
	stmt.d = i.getBits(1, 1)
	stmt.reg = 0b000
	switch stmt.w {
	case 0b0:
		i.Next()
		stmt.lo = (i.CurrentByte())
	case 0b1:
		i.Next()
		stmt.lo = (i.CurrentByte())
		i.Next()
		stmt.hi = (i.CurrentByte())
	default:
	}

	return &stmt
}
func (i *InstructionDecoder) MovIRegInstruction() InstructionStatement {
	stmt := MovInstruction{op: MOV_IMMEDIATE_TO_REGISTER}
	stmt.reg = i.getBits(0, 3)
	stmt.w = i.getBits(3, 1)

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
	stmt := MovInstruction{op: MOV_IMMEDIATE_TO_REGISTER_MEMORY}
	stmt.w = i.getBits(0, 1)
	// rm is the destination
	stmt.d = 0

	i.Next()
	stmt.mod = i.getBits(6, 2)

	stmt.rm = i.getBits(0, 3)

	switch stmt.mod {
	case 0b00:
		if stmt.rm == 0b110 {
			i.Next()
			stmt.lo = (i.CurrentByte())
			i.Next()
			stmt.hi = (i.CurrentByte())
		}
	case 0b01:
		i.Next()
		stmt.lo = (i.CurrentByte())
	case 0b10:
		i.Next()
		stmt.lo = (i.CurrentByte())
		i.Next()
		stmt.hi = (i.CurrentByte())
	}

	switch stmt.w {
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
func (i *InstructionDecoder) ArithRMAcc() InstructionStatement {
	v := i.getBits(2, 1)
	switch v {
	case 0b0:
		return i.ArithRMInstruction()
	case 0b1:
		return i.ArithAcc()
	}
	return nil
}

func (i *InstructionDecoder) ArithRMInstruction() InstructionStatement {
	stmt := ArithmeticInstruction{}
	stmt.w = i.getBits(0, 1)
	stmt.d = i.getBits(1, 1)
	OperationCode := i.getBits(3, 3)

	i.Next()
	stmt.mod = i.getBits(6, 2)
	stmt.reg = i.getBits(3, 3)
	stmt.rm = i.getBits(0, 3)
	switch stmt.mod {
	case 0b00:
		if stmt.rm == 0b110 {
			i.Next()
			stmt.lo = (i.CurrentByte())
			i.Next()
			stmt.hi = (i.CurrentByte())
		}
	case 0b01:
		i.Next()
		stmt.lo = (i.CurrentByte())
	case 0b10:
		i.Next()
		stmt.lo = (i.CurrentByte())
		i.Next()
		stmt.hi = (i.CurrentByte())
	}

	switch OperationCode {
	case 0b000:
		stmt.op = ADD_REG_MEM
	case 0b101:
		stmt.op = SUB_REG_MEM
	case 0b111:
		stmt.op = CMP_REG_MEM
	default:
		stmt.op = INSTRUCTION_UNKNOWN
	}

	return &stmt
}
func (i *InstructionDecoder) ArithImmediateInstruction() InstructionStatement {
	stmt := ArithmeticInstruction{}
	stmt.w = i.getBits(0, 1)
	stmt.s = i.getBits(1, 1)

	i.Next()
	stmt.mod = i.getBits(6, 2)
	stmt.rm = i.getBits(0, 3)
	OperationCode := i.getBits(3, 3)

	switch stmt.mod {
	case 0b00:
		if stmt.rm == 0b110 {
			i.Next()
			stmt.lo = (i.CurrentByte())
			i.Next()
			stmt.hi = (i.CurrentByte())
		}
	case 0b01:
		i.Next()
		stmt.lo = (i.CurrentByte())
	case 0b10:
		i.Next()
		stmt.lo = (i.CurrentByte())
		i.Next()
		stmt.hi = (i.CurrentByte())
	}

	switch stmt.w {
	case 0b1:
		if stmt.s == 1 {
			i.Next()
			stmt.data = int16(i.CurrentByte())
		} else {
			var b []uint16
			i.Next()
			b = append(b, uint16(i.CurrentByte()))
			i.Next()
			b = append(b, uint16(i.CurrentByte()))
			stmt.data = int16(b[1]<<8 | b[0])
		}
	case 0b0:
		i.Next()
		stmt.data = int16(i.CurrentByte())
	}

	switch OperationCode {
	case 0b000:
		stmt.op = ADD_IMMEDIATE_RM
	case 0b101:
		stmt.op = SUB_IMMEDIATE_RM
	case 0b111:
		stmt.op = CMP_IMMEDIATE_RM
	default:
		stmt.op = INSTRUCTION_UNKNOWN
	}
	return &stmt
}
func (i *InstructionDecoder) ArithAcc() InstructionStatement {
	stmt := ArithmeticInstruction{}
	stmt.w = i.getBits(0, 1)
	stmt.reg = 0b000
	OperationCode := i.getBits(3, 3)

	switch stmt.w {
	case 0b0:
		i.Next()
		stmt.lo = (i.CurrentByte())
	case 0b1:
		i.Next()
		stmt.lo = (i.CurrentByte())
		i.Next()
		stmt.hi = (i.CurrentByte())
	default:
	}

	switch OperationCode {
	case 0b000:
		stmt.op = ADD_ACC
	case 0b101:
		stmt.op = SUB_ACC
	case 0b111:
		stmt.op = CMP_ACC
	default:
		stmt.op = INSTRUCTION_UNKNOWN
	}
	return &stmt
}
func (i *InstructionDecoder) JumpLoop() InstructionStatement {
	stmt := JumpLoopInstruction{}
	stmt.op = i.getBits(0, 5)
	i.Next()
	stmt.ipInc = int8(i.CurrentByte())
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
	// fmt.Printf("%08b & %08b\n", i.CurrentByte()>>shift, (1<<nBits)-1)
	return (i.CurrentByte() >> shift) & ((1 << nBits) - 1)
}

func (i *InstructionDecoder) printCurrentByte() {
	fmt.Printf("%08b\n", i.CurrentByte())
}
