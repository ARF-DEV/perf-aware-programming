package internal

type InstructionStatement interface {
	String() string
	Disassemble() (string, error)
	Simulate(simulator *Simulator)
	isInstruction()
}

type Statements []InstructionStatement
