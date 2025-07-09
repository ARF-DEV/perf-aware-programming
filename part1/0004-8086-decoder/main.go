package main

import (
	"decoder/internal/v2"
	"os"
)

func main() {
	filename := os.Args[1]
	decoder := internal.NewDecoder(filename)

	decoder.Decode()
	if err := decoder.Disassemble(os.Stdout); err != nil {
		panic(err)
	}
}
