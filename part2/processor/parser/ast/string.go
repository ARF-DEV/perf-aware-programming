package ast

import (
	"fmt"
	"parttwo/processor/lexer"
)

type String struct {
	Token lexer.Token
	Value string
}

func (s *String) TokenValue() lexer.Value {
	return s.Token.Value
}

func (s *String) String() string {
	return fmt.Sprintf("%v %v", s.Token, s.Value)
}
