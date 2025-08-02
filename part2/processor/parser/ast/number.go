package ast

import (
	"fmt"
	"parttwo/processor/lexer"
)

type Integer struct {
	Token lexer.Token
	Value int64
}

func (s *Integer) String() string {
	return fmt.Sprintf("%v", s.Value)
}
func (n *Integer) TokenValue() lexer.Value {
	return n.Token.Value
}

// TODO: Parse Float, Parse Number (determined what parse function to use interger or float)
type Float struct {
	Token lexer.Token
	Value float64
}

func (n *Float) TokenValue() lexer.Value {
	return n.Token.Value
}

func (s *Float) String() string {
	return fmt.Sprintf("%v", s.Value)
}
