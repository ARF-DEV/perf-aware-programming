package ast

import (
	"fmt"
	"parttwo/processor/lexer"
)

type Integer struct {
	Token lexer.Token
	Value int64
}

func (n *Integer) TokenValue() lexer.Value {
	return n.Token.Value
}

func ParseIntegerNode(token *lexer.Token) (Node, error) {
	if len(token.Value) == 0 {
		return nil, fmt.Errorf("invalid number: %v", token.Value)
	}

	var literal string = string(token.Value)
	var isNegative bool
	if token.Value[0] == '-' {
		literal = literal[1:]
		isNegative = true
	}

	digits := len(literal)
	place := 1
	val := 0
	for i := 1; i <= digits; i++ {
		val += int(token.Value[digits-i]-'0') * place
		place *= 10
	}

	if isNegative {
		val *= -1
	}
	n := Integer{
		Token: *token,
		Value: int64(val),
	}
	return &n, nil
}

// TODO: Parse Float, Parse Number (determined what parse function to use interger or float)
