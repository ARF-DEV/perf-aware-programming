package internal

import (
	"fmt"
	"log"
)

var jumpOpMap map[byte]string = map[byte]string{
	// cond jumps
	0b10100: "je",
	0b11100: "jl",
	0b11110: "jle",
	0b10010: "jb",
	0b10110: "jbe",
	0b11010: "jp",
	0b10000: "jo",
	0b11000: "js",
	0b10101: "jne",
	0b11101: "jnl",
	0b11111: "jg",
	0b10011: "jnb",
	0b10111: "ja",
	0b11011: "jnp",
	0b10001: "jno",
	0b11001: "jns",

	// loops
	0b00010: "loop",
	0b00001: "loopz",
	0b00000: "loopnz",
	0b00011: "jcxz",
}

func (i *InstructionDecoder) Jump() string {
	subOp := i.getMaskedBits(0, 0b00011111)
	jumpOp, ok := jumpOpMap[subOp]
	if !ok {
		log.Fatalf("%08b not found", subOp)
	}
	i.Next()
	return fmt.Sprintf("%s %d", jumpOp, int8(i.CurrentByte()))
}

func (i *InstructionDecoder) Loop() string {
	subOp := i.getMaskedBits(0, 0b00011111)
	loopOp, ok := jumpOpMap[subOp]
	if !ok {
		log.Fatalf("%08b not found", subOp)
	}
	i.Next()
	return fmt.Sprintf("%s %d", loopOp, int8(i.CurrentByte()))
}
