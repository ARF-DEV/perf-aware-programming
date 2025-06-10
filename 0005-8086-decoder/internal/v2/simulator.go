package internal

import (
	"fmt"
	"strings"
)

type Simulator struct {
	register Register
	flags    Flags
	ip       *int64 // instruction pointer
	memory   Memory
}

func (s *Simulator) String() string {
	b := strings.Builder{}
	fmt.Fprintf(&b, "Final registers:\n%v", s.register)
	fmt.Fprintf(&b, "ip: 0x%04x (%d)\n", uint64(*s.ip), *s.ip)
	fmt.Fprintln(&b, "Final flags: ", s.flags)
	return b.String()
}

// func NewSimulator(statements Statements) InstructionSimulator {
// 	sim := InstructionSimulator{
// 		memory: Memory{},
// 		flags:  Flags{},
// 	}
// 	return sim
// }
// func (s *InstructionSimulator) Simulate() {
// 	fmt.Println(s.memory)
// 	for _, stmt := range s.statements {
// 		stmt.Simulate(&s.memory, &s.flags)
// 		// fmt.Println(stmt.String())
// 	}
// 	fmt.Println()
// 	fmt.Printf("Final registers:\n%v\n", s.memory)
// 	fmt.Println("Final flags: ", s.flags)
// }
