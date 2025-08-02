package main

import (
	"fmt"
	"log"
	"os"
	"parttwo/processor/lexer"
	"parttwo/processor/parser"
)

type A struct {
	Test int
	B    int
}

func main() {
	args := os.Args
	if len(args) < 2 {
		log.Fatal("cmd: main.go <file>")
	}
	input, err := os.ReadFile(args[1])
	if err != nil {
		log.Fatal(err)
	}

	l := lexer.New(string(input))
	l.Process()

	p := parser.New(&l)
	p.Process()
	var apa map[string]any = map[string]any{}
	if err := p.Decode(&apa); err != nil {
		log.Fatal(err)
	}

	fmt.Println(apa)
}
