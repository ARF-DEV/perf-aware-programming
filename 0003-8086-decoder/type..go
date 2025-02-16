package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type InstructionFunc func() string
type InstructionInfo struct {
	Length uint8
	Ins    InstructionFunc
}

// key = {RtoR}{w}{Reg/RM} -> [1bit][1bit][3bit]
var RegisterMap map[byte]string = map[byte]string{
	// MOD = 11
	0b10000: "al",
	0b10001: "cl",
	0b10010: "dl",
	0b10011: "bl",
	0b10100: "ah",
	0b10101: "ch",
	0b10110: "dh",
	0b10111: "bh",

	0b11000: "ax",
	0b11001: "cx",
	0b11010: "dx",
	0b11011: "bx",
	0b11100: "sp",
	0b11101: "bp",
	0b11110: "si",
	0b11111: "di",

	// the others i.e :
	// MOD = 00
	// MOD = 10
	// MOD = 01

	0b00000: "bx + si",
	0b00001: "bx + di",
	0b00010: "bp + si",
	0b00011: "bp + di",
	0b00100: "si",
	0b00101: "di",
	0b00110: "bp",
	0b00111: "bx",

	0b01000: "bx + si",
	0b01001: "bx + di",
	0b01010: "bp + si",
	0b01011: "bp + di",
	0b01100: "si",
	0b01101: "di",
	0b01110: "bp",
	0b01111: "bx",
}

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

	ins.initMap()
	ins.builder = strings.Builder{}
	ins.builder.WriteString("bits 16\n\n")
	return &ins
}

func (i *InstructionDecoder) initMap() {
	i.instructionFuncs = map[byte]InstructionFunc{
		0b100010:  i.MovInstruction,
		0b1011:    i.MovIRegInstruction,
		0b1100011: i.MovIRMInstruction,
		0b1010000: i.MovMToAcc,
		0b1010001: i.MovAccToM,

		0b000000:  i.IncRegToMem,
		0b001010:  i.IncRegToMem,
		0b001110:  i.IncRegToMem,
		0b100000:  i.ImmediateToRM,
		0b0000010: i.ImmediateToAcc,
		0b0010110: i.ImmediateToAcc,
		0b0011110: i.ImmediateToAcc,
	}
	i.instructionFuncss = map[int]map[byte]InstructionFunc{
		6: {
			0b100010: i.MovInstruction,
			0b000000: i.IncRegToMem,
			0b001010: i.IncRegToMem,
			0b001110: i.IncRegToMem,
			0b100000: i.ImmediateToRM,
		},
		4: {
			0b1011: i.MovIRegInstruction,
			0b0111: i.Jump,
			0b1110: i.Loop,
		},
		7: {
			0b1100011: i.MovIRMInstruction,
			0b1010000: i.MovMToAcc,
			0b1010001: i.MovAccToM,
			0b0000010: i.ImmediateToAcc,
			0b0010110: i.ImmediateToAcc,
			0b0011110: i.ImmediateToAcc,
		},
	}
}

func (i *InstructionDecoder) Decode() string {
	i.printBytes()
	for i.Next() {
		b := i.CurrentByte()
		for j := 0; j < 8; j++ {
			// ins, ok := i.instructionFuncs[b>>j]

			ins, ok := i.instructionFuncss[8-j][b>>j]
			if !ok {
				// fmt.Printf("conn %d, %b\n", 8-j, b>>j)
				continue
			} else {

				// fmt.Println(i.instructionFuncss[8-j], j)
			}
			res := ins()
			i.builder.WriteString(res + "\n")
			// i.builder.WriteString(res + fmt.Sprintf(" %d", i.curIdx) + "\n")
			break
		}
		// fmt.Println()
	}
	return i.builder.String()
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

func (i *InstructionDecoder) isDestination() bool {
	return (i.CurrentByte()>>1)&1 == 1
}

func (i *InstructionDecoder) isWide() bool {
	return i.getWide() == 1
}
func (i *InstructionDecoder) getWide() byte {
	return (i.CurrentByte() & 1)
}

func (i *InstructionDecoder) getReg() byte {
	return (i.CurrentByte() >> 3) & 0b00000111
}
func (i *InstructionDecoder) getRM() byte {
	return i.CurrentByte() & 0b00000111
}

func (i *InstructionDecoder) getMod() byte {
	return (i.CurrentByte() >> 6) & 0b00000011
}

func (i *InstructionDecoder) getMaskedBits(shift int64, mask byte) byte {
	return (i.CurrentByte() >> shift) & mask
}

func (i *InstructionDecoder) printBytes() {
	for i, v := range i.buf {
		fmt.Printf("%d. %08b ", i, v)
	}
	fmt.Println()
	// fmt.Printf("%08b\n", i.buf)
}
