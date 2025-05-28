package main

import (
	"decoder/internal/v2"
	"flag"
	"log"
	"os"
)

func main() {
	flagExec := flag.Bool("exec", false, "exec asm")
	flag.Parse()
	args := flag.Args()
	if len(args) != 1 {
		log.Fatalf("error: go run main.go <options> <file>\n")
		return
	}
	filename := flag.Args()[0]

	decoder := internal.NewDecoder(filename)
	decoder.Decode()
	if *flagExec {
		simulator := internal.NewSimulator(decoder.Statements())
		simulator.Simulate()
	} else {
		if err := decoder.Disassemble(os.Stdout); err != nil {
			panic(err)
		}
	}
}
