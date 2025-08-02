package ast

import (
	"parttwo/processor/lexer"
)

type KeyValuePair struct {
	token lexer.Token
	right Node
	left  Node
}
