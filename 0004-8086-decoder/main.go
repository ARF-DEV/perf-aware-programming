package main

import (
	"decoder/internal/v1"
	"fmt"
	"os"
)

func main() {
	filename := os.Args[1]
	decoder := internal.NewDecoder(filename)

	fmt.Print(decoder.Decode())
}
