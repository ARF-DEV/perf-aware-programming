package internal

import "fmt"

type InstructionSimulator struct {
	memory     Memory
	statements Statements
}

func NewSimulator(statements Statements) InstructionSimulator {
	sim := InstructionSimulator{
		memory:     Memory{},
		statements: statements,
	}
	return sim
}
func (s *InstructionSimulator) Simulate() {
	fmt.Println(s.memory)

	for _, stmt := range s.statements {
		// str, _ := stmt.Disassemble()
		stmt.Simulate(&s.memory)
		// fmt.Println(str)
	}

	fmt.Println(s.memory)
}
