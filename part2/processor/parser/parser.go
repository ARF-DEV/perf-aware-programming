package parser

import (
	"fmt"
	"log"
	"parttwo/processor/lexer"
	"parttwo/processor/parser/ast"
)

type Parser struct {
	lexer *lexer.Lexer
	Node  ast.Node

	idx      int
	curDepth int
}

func New(lexer *lexer.Lexer) *Parser {
	p := &Parser{
		lexer:    lexer,
		idx:      0,
		Node:     nil,
		curDepth: 0,
	}
	return p
}

func (p *Parser) Process() (err error) {
	parse, found := p.getNodeParseFunc(p.currentToken())
	if !found {
		return fmt.Errorf("error no opening '{'")
	}
	p.Node, err = parse()
	return err
}
func (p *Parser) parseObject() (ast.Node, error) {
	p.curDepth++
	node := ast.Object{LToken: *p.currentToken(), Depth: p.curDepth}
	p.nextIndex()
	for ; p.idx < len(p.lexer.Tokens); p.nextIndex() {
		p.skipToken(lexer.COMMA)
		if p.curTokenIs(lexer.RIGHT_CURLY_BRACKET) {
			break
		}
		parse, found := p.getNodeParseFunc(p.currentToken())
		if !found {
			log.Printf("debug: parser for token -> %v isn't implemented yet, depth: %d", p.currentToken(), p.curDepth)
			continue
		}
		newNode, err := parse()
		if err != nil {
			return nil, err
		}
		node.Values = append(node.Values, newNode)
	}
	node.RToken = *p.currentToken()
	p.curDepth--
	return &node, nil
}
func (p *Parser) parseArray() (ast.Node, error) {
	p.curDepth++
	node := ast.Array{LToken: *p.currentToken(), Depth: p.curDepth}
	p.nextIndex()
	for ; p.idx < len(p.lexer.Tokens); p.nextIndex() {
		p.skipToken(lexer.COMMA)
		if p.curTokenIs(lexer.RIGHT_SQUARE_BRACKET) {
			break
		}
		parse, found := p.getNodeParseFunc(p.currentToken())
		if !found {
			log.Printf("debug: parser for token -> %v isn't implemented yet, depth: %d", p.currentToken(), p.curDepth)
			continue
		}
		newNode, err := parse()
		if err != nil {
			return nil, err
		}
		node.Values = append(node.Values, newNode)
	}
	node.RToken = *p.currentToken()
	p.curDepth--
	return &node, nil
}

func (p *Parser) skipToken(tokenType lexer.Type) {
	for p.currentToken().Type == tokenType {
		p.nextIndex()
	}
}
func (p *Parser) curTokenIs(tokenType lexer.Type) bool {
	return p.currentToken().Type == tokenType
}

func (p *Parser) currentToken() *lexer.Token {
	return &p.lexer.Tokens[p.idx]
}

func (p *Parser) nextToken() *lexer.Token {
	if p.idx >= len(p.lexer.Tokens) {
		return nil
	}
	return &p.lexer.Tokens[p.idx+1]
}
func (p *Parser) nextIndex() {
	p.idx++
}

func (p *Parser) getNodeParseFunc(token *lexer.Token) (ast.NodeParseFunc, bool) {
	var nodeParseFuncMap map[lexer.Type]ast.NodeParseFunc = map[lexer.Type]ast.NodeParseFunc{
		lexer.STRING:              p.parseStringNode,
		lexer.NUMBER:              p.parseNumberNode,
		lexer.LEFT_CURLY_BRACKET:  p.parseObject,
		lexer.LEFT_SQUARE_BRACKET: p.parseArray,
	}
	parse, found := nodeParseFuncMap[token.Type]
	return parse, found
}
func (p *Parser) parseKeyValuePair(left ast.Node) (ast.Node, error) {
	// parse prev token, and next token
	node := ast.KeyValuePair{
		Token: *p.currentToken(),
		Left:  left,
	}

	p.nextIndex()

	parser, found := p.getNodeParseFunc(p.currentToken())
	if found {
		right, err := parser()
		if err != nil {
			return nil, err
		}
		node.Right = right
	}

	return &node, nil
}
func (p *Parser) parseNumberNode() (ast.Node, error) {
	token := p.currentToken()
	containsDot := false
	for _, c := range token.Value {
		if c == '.' {
			containsDot = true
			break
		}
	}
	if containsDot {
		return p.parseFloatNode()
	}

	return p.parseIntegerNode()
}
func (p *Parser) parseIntegerNode() (ast.Node, error) {
	token := p.currentToken()
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
	n := &ast.Integer{
		Token: *token,
		Value: int64(val),
	}
	if p.nextToken().Type == lexer.COLON {
		p.nextIndex()
		return p.parseKeyValuePair(n)
	}
	return n, nil
}

func (p *Parser) parseStringNode() (ast.Node, error) {
	token := p.currentToken()
	s := &ast.String{
		Token: *token,
		Value: string(token.Value),
	}

	if p.nextToken().Type == lexer.COLON {
		p.nextIndex()
		return p.parseKeyValuePair(s)
	}
	return s, nil
}

func (p *Parser) parseFloatNode() (ast.Node, error) {
	token := p.currentToken()
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
	node := &ast.Float{Token: *token, Value: float64(val) / float64(decimal)}
	if p.nextToken().Type == lexer.COLON {
		p.nextIndex()
		return p.parseKeyValuePair(node)
	}
	return node, nil
}
