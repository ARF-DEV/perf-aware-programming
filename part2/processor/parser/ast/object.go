package ast

import (
	"fmt"
	"parttwo/processor/lexer"
	"strings"
)

type Object struct {
	LToken, RToken lexer.Token
	Values         Nodes
	Depth          int
}

func (s *Object) TokenValue() lexer.Value {
	return s.LToken.Value
}

func (s *Object) String() string {
	if len(s.Values) == 0 {
		return "{}"
	}
	var str string = fmt.Sprintf("%v\n", s.LToken.Value)
	tabs := strings.Repeat("\t", s.Depth)
	for _, v := range s.Values {
		str += fmt.Sprintf("%s%v\n", tabs, v.String())
	}
	str += fmt.Sprint(strings.Repeat("\t", s.Depth-1), s.RToken.Value)

	return str
}

type Array struct {
	LToken, RToken lexer.Token
	Values         Nodes
	Depth          int
}

func (s *Array) TokenValue() lexer.Value {
	return s.LToken.Value
}

func (s *Array) String() string {
	if len(s.Values) == 0 {
		return "[]"
	}
	var str string = fmt.Sprintf("%v\n", s.LToken.Value)
	tabs := strings.Repeat("\t", s.Depth)
	for _, v := range s.Values {
		str += fmt.Sprintf("%s%v\n", tabs, v.String())
	}
	str += fmt.Sprint(strings.Repeat("\t", s.Depth-1), s.RToken.Value)

	return str
}
