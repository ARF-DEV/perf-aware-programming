package internal

type InstructionStatement interface {
	String() string
	Disassemble() (string, error)
	Simulate(Memory *Memory)
	isInstruction()
}

type Statements []InstructionStatement
