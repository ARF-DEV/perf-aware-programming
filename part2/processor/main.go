package main

import (
	"fmt"
	"log"
	"os"
	"parttwo/processor/lexer"
)

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
	fmt.Println(l.Tokens)
}
