package ast

import (
	"fmt"
	"parttwo/processor/lexer"
)

func ParseNumberNode(token *lexer.Token) (Node, error) {
	containsDot := false
	for _, c := range token.Value {
		if c == '.' {
			containsDot = true
			break
		}
	}
	if containsDot {
		return parseFloatNode(token)
	}

	return parseIntegerNode(token)
}

type Integer struct {
	Token lexer.Token
	Value int64
}

func (n *Integer) TokenValue() lexer.Value {
	return n.Token.Value
}

func parseIntegerNode(token *lexer.Token) (Node, error) {
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
		c := literal[digits-i]
		if c < '0' || c > '9' {
			return nil, fmt.Errorf("invalid number: %v", token.Value)
		}
		val += int(c-'0') * place
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
type Float struct {
	Token lexer.Token
	Value float64
}

func (n *Float) TokenValue() lexer.Value {
	return n.Token.Value
}

func parseFloatNode(token *lexer.Token) (Node, error) {
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
		c := literal[digits-i]
		if c == '.' {
			continue
		}
		if c < '0' || c > '9' {
			return nil, fmt.Errorf("invalid number: %v", token.Value)
		}
		val += int(c-'0') * place
		place *= 10
	}
	if isNegative {
		val *= -1
	}

	decimal := 1
	for i := digits - 1; literal[i] != '.'; i-- {
		decimal *= 10
	}

	return &Float{Token: *token, Value: float64(val) / float64(decimal)}, nil
}
