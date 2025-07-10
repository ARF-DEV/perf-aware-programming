package main

import (
	"fmt"
	"parttwo/processor/lexer"
)

func main() {
	//	input := `
	//
	// [
	//
	//	{
	//		"point_1": {
	//			"latitude": 20.08928198951905,
	//			"longitude": -175.29980187201667
	//		},
	//		"point_2": {
	//			"latitude": 72.72640880650494,
	//			"longitude": -65.88179901773461
	//		}
	//	}
	//
	// ]
	//
	//	`
	// input := "{ } [ ]"
	// input := `{"key": "value"}`
	// input := `{"key": 121231 }`
	input := `
		{
  "key1": true,
  "key2": false,
  "key3": null,
  "key4": "value",
  "key5": 101
}	
	`
	fmt.Println(input)
	l := lexer.New(input)
	l.Process()
	fmt.Println(l.Tokens)
}
