package main

import "fmt"

type InstructionFunc func([2]byte) string

var InstructionMap map[byte]InstructionFunc

// key = {w}{register} -> [1bit][3bit]
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

func init() {
	InstructionMap = make(map[byte]InstructionFunc, 0)
}

func MovInstruction(instructionByte [2]byte) string {
	Reg := getReg(instructionByte[1])
	RM := getRM(instructionByte[1])
	w := getWide(instructionByte[0])
	d := isDestination(instructionByte[0])

	RegStr := RegisterMap[Reg|(w<<3)]
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

func isDestination(b byte) bool {
	return b&(1<<1) == 1
}

func isWide(b byte) bool {
	return getWide(b) == 1
}
func getWide(b byte) byte {
	return (b & 1)
}

func getReg(b byte) byte {
	return (b >> 3) & 0b00000111
}
func getRM(b byte) byte {
	return b & 0b00000111
}

// func getMode() {

// }
