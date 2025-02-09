package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type InstructionFunc func([2]byte) string

var InstructionMap map[byte]InstructionFunc

// key = {w}{register} -> [1bit][3bit]
// TBD: key = {MOD}{w}{Reg/RM} -> [2bit][1bit][3bit]
var RegisterMap map[byte]string = map[byte]string{
	0b0000: "al",
	0b0001: "cl",
	0b0010: "dl",
	0b0011: "bl",
	0b0100: "ah",
	0b0101: "ch",
	0b0110: "dh",
	0b0111: "bh",

	0b1000: "ax",
	0b1001: "cx",
	0b1010: "dx",
	0b1011: "bx",
	0b1100: "sp",
	0b1101: "bp",
	0b1110: "si",
	0b1111: "di",
}

type InstructionDecoder struct {
	buf     []byte
	curIdx  int64
	nextIdx int64
	builder strings.Builder
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

func (i *InstructionDecoder) Decode() string {
	for i.Next() {
		res := i.MovInstruction()
		i.builder.WriteString(res + "\n")
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

func init() {
	InstructionMap = make(map[byte]InstructionFunc, 0)
}

func (i *InstructionDecoder) MovInstruction() string {
	w := i.getWide()
	d := i.isDestination()

	i.Next()
	Reg := i.getReg()
	RegStr := RegisterMap[Reg|(w<<3)]

	RM := i.getRM()
	RMStr := RegisterMap[RM|(w<<3)]

	var dst, src string
	if d {
		dst = RegStr
		src = RMStr
	} else {
		src = RegStr
		dst = RMStr
	}

	return fmt.Sprintf("mov %s, %s", dst, src)
}

func (i *InstructionDecoder) isDestination() bool {
	return i.CurrentByte()&(1<<1) == 1
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
