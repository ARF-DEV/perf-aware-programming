package internal

type InstructionFunc func()
type OpMode uint8

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
