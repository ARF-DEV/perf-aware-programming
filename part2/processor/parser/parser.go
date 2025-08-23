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
	switch rv.Kind() {
	case reflect.Map:
		return parseToMap(p.Node, rv)
	case reflect.Slice:
		return parseToArray(p.Node, rv)
	case reflect.Struct:
		// TODO
		return parseToStruct(p.Node, rv)
	}

	return fmt.Errorf("cannot process value of type %v", rv.Type())
}

func parseToMap(node ast.Node, m reflect.Value) error {
	switch actual := node.(type) {
	case *ast.Object:
		err := parseObject(actual, m)
		if err != nil {
			m.Set(reflect.MakeMap(m.Type()))
			return err
		}
		return nil
	default:
		return fmt.Errorf("cannot parse node of type %T to %v", node, m.Type())
	}
}
func parseToArray(node ast.Node, s reflect.Value) error {
	switch actual := node.(type) {
	case *ast.Array:
		values, err := parseArray(actual)
		if err != nil {
			return err
		}
		s.Set(reflect.ValueOf(values))
	default:
		return fmt.Errorf("cannot parse node of type %T to %v", node, s.Type())
	}
	return nil
}
func parseToStruct(node ast.Node, st reflect.Value) error {
	switch actual := node.(type) {
	case *ast.Object:
		return parseStruct(actual, st)
	default:
		return fmt.Errorf("cannot parse node of type %T to %v", node, st.Type())
	}
}
func parseStruct(obj *ast.Object, st reflect.Value) error {
	nodeMap := map[lexer.Value]*ast.KeyValuePair{}
	for _, kv := range obj.Values {
		actualKv, isPair := kv.(*ast.KeyValuePair)
		if !isPair {
			log.Printf("%v is not key value pair", kv)
			continue
		}
		nodeMap[actualKv.Left.TokenValue()] = actualKv
	}

	len := st.NumField()
	stv := st.Type()
	for i := 0; i < len; i++ {
		field := st.Field(i)
		fieldType := stv.Field(i)

		// fmt.Printf("%+v: %+v\n", fieldType.Tag.Get("json"), nodeMap[lexer.Value(fieldType.Tag.Get("json"))])
		key := lexer.Value(fieldType.Tag.Get("json"))
		kvValue, found := nodeMap[key]
		if !found {
			log.Printf("key %v not found", key)
			continue
		}
		var value any
		switch tkActual := kvValue.Right.(type) {
		case *ast.String:
			value = tkActual.Value
		case *ast.Float:
			value = tkActual.Value
		case *ast.Integer:
			value = tkActual.Value
		case *ast.Object:
			parseStruct(tkActual, field)
			continue
		case *ast.Array:
			v, err := parseArray(tkActual)
			if err != nil {
				return err
			}
			value = v
		}

		if err := setField(field, reflect.ValueOf(value)); err != nil {
			return err
		}

	}

	// fmt.Printf("final %+v\n", st)
	return nil
}

func parseObject(obj *ast.Object, m reflect.Value) (err error) {
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
		case *ast.Object:
			childMap := map[string]any{}
			rcm := reflect.ValueOf(&childMap).Elem()
			if err := parseToMap(actualVal, rcm); err != nil {
				return err
			}
			value = childMap
		case *ast.Array:
			v, err := parseArray(actualVal)
			if err != nil {
				return err
			}
			value = v

		default:
			log.Printf("debug: not implemented for %v", actualVal)
			continue
		}
		if err := setMap(m, reflect.ValueOf(key.Value), reflect.ValueOf(value)); err != nil {
			return err
		}
	}
	return nil
}

func parseArray(arrayNode *ast.Array) ([]any, error) {
	values := []any{}
	for _, arrayValue := range arrayNode.Values {
		var value any
		switch actualVal := arrayValue.(type) {
		case *ast.String:
			value = actualVal.Value
		case *ast.Float:
			value = actualVal.Value
		case *ast.Integer:
			value = actualVal.Value
		case *ast.Object:
			childMap := map[string]any{}
			rcm := reflect.ValueOf(childMap)
			if err := parseToMap(actualVal, rcm); err != nil {
				return nil, err
			}
			value = childMap
		// case *ast.Array:
		default:
			log.Printf("debug: not implemented for %v", actualVal)
			continue
		}

		values = append(values, value)
	}
	return values, nil
}

func setField(f reflect.Value, v reflect.Value) error {
	if !v.Type().AssignableTo(f.Type()) && !v.Type().ConvertibleTo(f.Type()) {
		return fmt.Errorf("value of type %v cannot be assign to field of type %v", v.Type(), f.Type())
	}
	switch f.Kind() {
	case reflect.String:
		f.SetString(v.String())
	case reflect.Float32, reflect.Float64:
		kind := v.Kind()
		if kind != reflect.Float32 && kind != reflect.Float64 {
			v = intToFloat(v)
		}
		f.SetFloat(v.Float())
	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int8:
		kind := v.Kind()
		if kind != reflect.Int && kind != reflect.Int8 && kind != reflect.Int16 &&
			kind != reflect.Int32 && kind != reflect.Int64 {
			v = floatToInt(v)
		}
		f.SetInt(v.Int())
	}
	return nil
}

func floatToInt(v reflect.Value) reflect.Value {
	value := int64(v.Float())
	return reflect.ValueOf(value)
}
func intToFloat(v reflect.Value) reflect.Value {
	value := float64(v.Int())
	return reflect.ValueOf(value)
}

func setMap(m reflect.Value, k, v reflect.Value) error {
	if !k.Type().AssignableTo(m.Type().Key()) {
		return fmt.Errorf("value of type %v cannot be assign to field of type %v", k.Type(), m.Type().Key())
	}
	if !v.Type().AssignableTo(m.Type().Elem()) && !v.Type().ConvertibleTo(m.Type().Elem()) {
		return fmt.Errorf("value of type %v cannot be assign to field of type %v", v.Type(), m.Type().Elem())
	}
	target := m.Type().Elem()
	switch target.Kind() {
	case reflect.String, reflect.Slice, reflect.Array, reflect.Struct:
	case reflect.Float32, reflect.Float64:
		kind := v.Kind()
		if kind != reflect.Float32 && kind != reflect.Float64 {
			v = intToFloat(v)
		}
	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int8:
		kind := v.Kind()
		if kind != reflect.Int && kind != reflect.Int8 && kind != reflect.Int16 &&
			kind != reflect.Int32 && kind != reflect.Int64 {
			v = floatToInt(v)
		}
	}
	m.SetMapIndex(k, v)
	return nil
}
