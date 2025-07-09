package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	f, err := os.Open("listing1")
	if err != nil {
		panic(err)
	}
	info, err := f.Stat()
	if err != nil {
		panic(err)
	}
	r := bufio.NewReader(f)
	buf := make([]byte, info.Size())
	_, err = r.Read(buf)
	if err != nil {
		panic(err)
	}

	results := []string{
		"bits 16\n\n",
	}
	for i := 0; i < len(buf); i += 2 {
		high, low := buf[i], buf[i+1]
		res := MovInstruction([2]byte{high, low})
		results = append(results, res)
	}
	fmt.Print(strings.Join(results, "\n"))
}
