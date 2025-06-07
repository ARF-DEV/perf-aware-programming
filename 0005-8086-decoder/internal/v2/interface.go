package internal

type InstructionStatement interface {
	String() string
	Disassemble() (string, error)
	Simulate(Memory *Memory, flags *Flags)
	isInstruction()
}

type Statements []InstructionStatement
