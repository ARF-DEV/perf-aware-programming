package internal

type InstructionStatement interface {
	String() string
	Disassemble(clocks *int) (string, error)
	Simulate(simulator *Simulator)
	isInstruction()
}

type Statements []InstructionStatement
