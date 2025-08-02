package parser

import (
	"log"
	"parttwo/processor/lexer"
	"parttwo/processor/parser/ast"
)

type Parser struct {
	lexer *lexer.Lexer
	Nodes ast.Nodes

	idx int
}

func New(lexer *lexer.Lexer) *Parser {
	p := &Parser{
		lexer: lexer,
		Nodes: ast.Nodes{},
		idx:   0,
	}
	return p
}

func (p *Parser) Process() error {
	for ; p.idx < len(p.lexer.Tokens); p.nextIndex() {
		parse, found := p.getNodeParseFunc(p.currentToken())
		if !found {
			log.Printf("debug: parser for token -> %v isn't implemented yet", p.currentToken())
			continue
		}
		newNode, err := parse(p.currentToken())
		if err != nil {
			return err
		}
		p.Nodes = append(p.Nodes, newNode)
	}
	return nil
}

func (p *Parser) currentToken() *lexer.Token {
	return &p.lexer.Tokens[p.idx]
}
func (p *Parser) nextIndex() {
	p.idx++
}

func (p *Parser) getNodeParseFunc(token *lexer.Token) (ast.NodeParseFunc, bool) {
	parse, found := nodeParseFuncMap[token.Type]
	return parse, found
}

var nodeParseFuncMap map[lexer.Type]ast.NodeParseFunc = map[lexer.Type]ast.NodeParseFunc{
	lexer.STRING: ast.ParseStringNode,
	lexer.NUMBER: ast.ParseNumberNode,
}
