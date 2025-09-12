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
	Scores []int `json:"scores"`
}

type Scores struct {
	Values []float64 `json:"values"`
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
	// apa := []int{} // error
	// apa := []float64{} // integer to float and the reverse
	// apa := map[string]any{}
	apa := Person{}
	// apa := Scores{}
	if err := p.Decode(&apa); err != nil {
		log.Println(err)
	}
	fmt.Printf("%+v\n", apa)

}
