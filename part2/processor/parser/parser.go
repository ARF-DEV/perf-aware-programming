package parser

import (
	"fmt"
	"log"
	"parttwo/processor/lexer"
	"parttwo/processor/parser/ast"
	"reflect"
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

func (p *Parser) Decode(v any) error {
	if p.Node == nil {
		return fmt.Errorf("data empty")
	}
	rv := reflect.ValueOf(v)
	vt := rv.Type()
	if vt.Kind() != reflect.Pointer {
		return fmt.Errorf("the argument passed is a copy value")
	}

	rv = rv.Elem()
	// TODO: parse to struct and array
	return p.parseToMap(rv)
}

func (p *Parser) parseToMap(m reflect.Value) error {
	obj, ok := p.Node.(*ast.Object)
	if !ok {
		return fmt.Errorf("debug: not an object")
	}

	for _, value := range obj.Values {
		kp, ok := value.(*ast.KeyValuePair)
		if !ok {
			return fmt.Errorf("debug: not a key value pair")
		}

		key, ok := kp.Left.(*ast.String)
		if !ok {
			return fmt.Errorf("debug: not a string")
		}

		var value any
		switch actualVal := kp.Right.(type) {
		case *ast.String:
			value = actualVal.Value
		case *ast.Float:
			value = actualVal.Value
		case *ast.Integer:
			value = actualVal.Value
			// TODO: object and array
		default:
			log.Printf("debug: not implemented for %v", actualVal)
			continue
		}
		fmt.Println(key.Value, value)
		if err := p.setMap(m, reflect.ValueOf(key.Value), reflect.ValueOf(value)); err != nil {
			return err
		}
	}
	return nil
}

func (p *Parser) setMap(m reflect.Value, k, v reflect.Value) error {
	if !k.Type().AssignableTo(m.Type().Key()) {
		return fmt.Errorf("value of type %v cannot be assign to field of type %v", k.Type(), m.Type().Key())
	}
	if !v.Type().AssignableTo(m.Type().Elem()) {
		return fmt.Errorf("value of type %v cannot be assign to field of type %v", v.Type(), m.Type().Elem())
	}

	m.SetMapIndex(k, v)
	return nil
}
