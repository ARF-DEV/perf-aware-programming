package main

import (
	"fmt"
	"log"
	"os"
	"parttwo/processor/lexer"
	"parttwo/processor/parser"
)

type Person struct {
	Name      string  `json:"name"`
	Age       int     `json:"age"`
	Balance   float64 `json:"balance"`
	Education struct {
		InstitutionName string `json:"institution_name"`
		Degree          string `json:"degree"`
	} `json:"current_education"`
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
	// var apa map[string]float64 = map[string]float64{}
	// apa := []any{}
	apa := Person{}
	if err := p.Decode(&apa); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", apa)
}
