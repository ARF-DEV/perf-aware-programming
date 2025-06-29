package internal

import (
	"fmt"
	"log"
)

var jumpLoopTable map[byte]string = map[byte]string{
	0b10100: "je",
	0b11100: "jl",
	0b11110: "jng",
	0b10010: "jb",
	0b10110: "jna",
	0b11010: "jp",
	0b10000: "jo",
	0b11000: "js",
	0b10101: "jnz",
	0b11101: "jge",
	0b11111: "jg",
	0b10011: "jae",
	0b10111: "ja",
	0b11011: "jpo",
	0b10001: "jno",
	0b11001: "jns",
	0b00010: "loop",
	0b00001: "loopz",
	0b00000: "loopnz",
	0b00011: "jcxz",
}

func (i *JumpLoopInstruction) jumpLoopCondTable() map[byte]func(f *Flags) bool {
	return map[byte]func(f *Flags) bool{
		// 0b10100: "je",
		// 0b11100: "jl",
		// 0b11110: "jng",
		// 0b10010: "jb",
		// 0b10110: "jna",
		// 0b11010: "jp",
		// 0b10000: "jo",
		// 0b11000: "js",
		0b10101: func(f *Flags) bool { // not zero
			return !f.Get(FLAGS_ZERO)
		},
		// 0b11101: "jge",
		// 0b11111: "jg",
		// 0b10011: "jae",
		// 0b10111: "ja",
		// 0b11011: "jpo",
		// 0b10001: "jno",
		// 0b11001: "jns",
		// 0b00010: "loop",
		// 0b00001: "loopz",
		// 0b00000: "loopnz",
		// 0b00011: "jcxz",
	}
}

type JumpLoopInstruction struct {
	opByte byte
	op     OpMode
	ipInc  int8
}

func (i *JumpLoopInstruction) String() string {
	return ""
}

func (i *JumpLoopInstruction) Disassemble(clocks *int) (string, error) {
	opStr := jumpLoopTable[i.opByte]
	return fmt.Sprintf("%s %d", opStr, i.ipInc), nil
}

func (i *JumpLoopInstruction) isInstruction() {}
func (i *JumpLoopInstruction) Simulate(simulator *Simulator) {
	opStr := jumpLoopTable[i.opByte]
	cond, found := i.jumpLoopCondTable()[i.opByte]
	if !found {
		log.Printf("jump %s not found\n", opStr)
		return
	}

	fmt.Printf("%s %d; ", opStr, i.ipInc)

	if !cond(&simulator.flags) {
		return
	}
	*simulator.ip += int64(i.ipInc)
}
