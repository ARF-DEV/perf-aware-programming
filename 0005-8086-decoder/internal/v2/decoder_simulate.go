package internal

import "fmt"

type InstructionSimulator struct {
	memory     Memory
	flags      Flags
	statements Statements
}

func NewSimulator(statements Statements) InstructionSimulator {
	sim := InstructionSimulator{
		memory:     Memory{},
		statements: statements,
		flags:      Flags{},
	}
	return sim
}
func (s *InstructionSimulator) Simulate() {
	fmt.Println(s.memory)
	for _, stmt := range s.statements {
		stmt.Simulate(&s.memory, &s.flags)
		// fmt.Println(stmt.String())
	}
	fmt.Println()
	fmt.Printf("Final registers:\n%v\n", s.memory)
	fmt.Println("Final flags: ", s.flags)
}
