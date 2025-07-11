package main

import (
	"decoder/internal/v2"
	"flag"
	"log"
	"os"
)

func main() {
	flagExec := flag.Bool("exec", false, "exec asm")
	flagDump := flag.Bool("dump", false, "dump final memory state")
	flag.Parse()
	args := flag.Args()
	if len(args) != 1 {
		log.Fatalf("error: go run main.go <options> <file>\n")
		return
	}
	filename := flag.Args()[0]

	decoder := internal.NewDecoder(filename)
	decoder.Decode(*flagExec, *flagDump)
	if *flagExec {
		return
	}
	if err := decoder.Disassemble(os.Stdout); err != nil {
		panic(err)
	}
}
