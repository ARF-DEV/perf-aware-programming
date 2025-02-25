package internal

type InstructionStatement interface {
	String() string
	Disassemble() (string, error)
	isInstruction()
}
