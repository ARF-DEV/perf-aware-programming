package internal

import (
	"fmt"
	"strings"
)

type InstructionFunc func() InstructionStatement
type OpMode uint8
type RegisterTable map[byte]string

type DecoderFunc func() string
type DecoderFuncTable map[OpMode]DecoderFunc

type SimulateFunc func(mem *Memory, flags *Flags)
type SimulateFuncTable map[OpMode]SimulateFunc

func (t RegisterTable) Get(mod, w, reg byte) string {
	var noDisp byte
	if mod == 0b11 {
		noDisp = 1
	} else {
		noDisp = 0
	}
	return t[noDisp<<4|w<<3|reg]
}

// key = {RtoR}{w}{Reg/RM} -> [1bit][1bit][3bit]
var RegisterTab RegisterTable = map[byte]string{
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

const (
	INSTRUCTION_UNKNOWN OpMode = iota
	MOV_REGISTER_FROM_TO_MEMORY
	MOV_IMMEDIATE_TO_REGISTER
	MOV_IMMEDIATE_TO_REGISTER_MEMORY
	MOV_ACCULUMATOR_FROM_TO_MEMORY

	ADD_REG_MEM
	ADD_IMMEDIATE_RM
	ADD_ACC
	SUB_REG_MEM
	SUB_IMMEDIATE_RM
	SUB_ACC
	CMP_REG_MEM
	CMP_IMMEDIATE_RM
	CMP_ACC
)

type Memory [8]int16
type Flags [2]bool
type FlagsIdx int

func (f Flags) Get(idx FlagsIdx) bool {
	return f[idx]
}

func (f *Flags) Set(idx FlagsIdx, val bool) {
	f[idx] = val
}
func (f Flags) String() string {
	str := ""
	for i, v := range f {
		if !v {
			continue
		}
		str += FLAGS_NAME[i]
	}
	return str
}

const (
	FLAGS_ZERO FlagsIdx = iota
	FLAGS_SIGN
)

var FLAGS_NAME = [2]string{"Z", "S"}
var REGISTERS_NAME = [8]string{"ax", "cx", "dx", "bx", "sp", "bp", "si", "di"}

func (m Memory) String() string {

	// Final registers:
	//   ax: 0x0001 (1)
	//   bx: 0x0002 (2)
	//   cx: 0x0003 (3)
	//   dx: 0x0004 (4)
	//   sp: 0x0005 (5)
	//   bp: 0x0006 (6)
	//   si: 0x0007 (7)
	//   di: 0x0008 (8)
	builder := strings.Builder{}

	for i := range m {
		fmt.Fprintf(&builder, "%s: 0x%04x (%d)\n", REGISTERS_NAME[i], m[i], m[i])
	}

	return builder.String()
}
