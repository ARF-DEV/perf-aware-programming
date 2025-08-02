package ast

import (
	"fmt"
	"parttwo/processor/lexer"
)

type KeyValuePair struct {
	Token lexer.Token
	Right Node
	Left  Node
}

func (n *KeyValuePair) TokenValue() lexer.Value {
	return n.Token.Value
}

func (s *KeyValuePair) String() string {
	return fmt.Sprintf("%v  %v  %v", s.Left, s.Token.Value, s.Right)
}
