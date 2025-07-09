package main

import (
	"fmt"
	"os"
)

func main() {
	filename := os.Args[1]
	decoder := NewDecoder(filename)

	fmt.Print(decoder.Decode())
}
