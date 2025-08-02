package ast

import (
	"fmt"
	"parttwo/processor/lexer"
)

type Node interface {
	TokenValue() lexer.Value
}
type Nodes []Node

func (n Nodes) String() string {
	str := ""
	for _, node := range n {
		str += fmt.Sprintln(node)
	}
	return str
}

type NodeParseFunc func(*lexer.Token) (Node, error)
