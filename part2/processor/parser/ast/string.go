package ast

import "parttwo/processor/lexer"

type String struct {
	Token *lexer.Token
	Value string
}

func (s *String) TokenValue() lexer.Value {
	return s.Token.Value
}

func ParseStringNode(token *lexer.Token) (Node, error) {
	s := &String{
		Token: token,
		Value: string(token.Value),
	}
	return s, nil
}
