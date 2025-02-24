package internal

type InstructionStatement interface {
	String() string
	isInstruction()
}
