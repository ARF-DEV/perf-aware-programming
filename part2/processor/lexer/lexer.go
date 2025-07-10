package lexer

type Lexer struct {
	Input  string
	Tokens Tokens

	idx int
}

func New(input string) Lexer {
	l := Lexer{
		Input:  input,
		Tokens: Tokens{},
		idx:    0,
	}
	return l
}

func (l *Lexer) Process() {
	for {
		l.skipWhiteSpace()
		if l.idx >= len(l.Input) {
			break
		}
		var newToken Token
		switch l.currentChar() {
		case '{':
			newToken = Token{
				Type:  LEFT_CURLY_BRACKET,
				Value: Value(l.currentChar()),
			}
		case '}':
			newToken = Token{
				Type:  RIGHT_CURLY_BRACKET,
				Value: Value(l.currentChar()),
			}
		case '[':
			newToken = Token{
				Type:  LEFT_SQUARE_BRACKET,
				Value: Value(l.currentChar()),
			}
		case ']':
			newToken = Token{
				Type:  RIGHT_SQUARE_BRACKET,
				Value: Value(l.currentChar()),
			}
		case ',':
			newToken = Token{
				Type:  COMMA,
				Value: Value(l.currentChar()),
			}
		case ':':
			newToken = Token{
				Type:  COLON,
				Value: Value(l.currentChar()),
			}
		default:
			newToken = l.getToken()
		}
		l.Tokens = append(l.Tokens, newToken)
		l.increment()
	}
}

func (l *Lexer) skipWhiteSpace() {
	for ; l.idx < len(l.Input); l.increment() {
		if !l.isWhiteSpace(l.currentChar()) {
			break
		}

	}
}

func (l *Lexer) getToken() Token {
	if l.isNumber(l.currentChar()) {
		return Token{
			Type:  NUMBER,
			Value: Value(l.getNumber()),
		}
	} else if l.isString(l.currentChar()) {
		return Token{
			Type:  STRING,
			Value: Value(l.getString()),
		}
	}
	val := l.getLiteral()
	switch val {
	case "true":
		return Token{
			Type:  TRUE,
			Value: Value(val),
		}
	case "false":
		return Token{
			Type:  FALSE,
			Value: Value(val),
		}
	case "null":
		return Token{
			Type:  NULL,
			Value: Value(val),
		}
	}
	return Token{
		Type:  UNKNOWN,
		Value: Value(val),
	}
}
func (l *Lexer) getNumber() string {
	number := ""
	for ; l.idx < len(l.Input); l.increment() {

		number += string(l.currentChar())

		if !l.isNumber(l.nextChar()) {
			break
		}

		if l.isComma(l.nextChar()) {
			break
		}
	}
	return number
}
func (l *Lexer) getLiteral() string {
	lit := ""
	for ; l.idx < len(l.Input); l.increment() {
		lit += string(l.currentChar())

		if l.isWhiteSpace(l.nextChar()) {
			break
		}
		if l.isComma(l.nextChar()) {
			break
		}
	}

	return lit
}
func (l *Lexer) getString() string {
	quote := string(l.currentChar())
	l.increment()
	str := ""
	for ; l.idx < len(l.Input) && quote != string(l.currentChar()); l.increment() {
		str += string(l.currentChar())
	}
	return str
}
func (l *Lexer) isNumber(c byte) bool {
	return (c >= '0' && c <= '9') || c == '-' || c == '.'
}
func (l *Lexer) isString(c byte) bool {
	return c == '"' || c == '\''
}

func (l *Lexer) isComma(c byte) bool {
	return c == ','
}
func (l *Lexer) isWhiteSpace(c byte) bool {
	return c == ' ' ||
		c == '\n' ||
		c == '\t' ||
		c == '\r' ||
		c == '\f' ||
		c == '\v'
}

func (l *Lexer) increment() {
	l.idx++
}

func (l *Lexer) currentChar() byte {
	return l.Input[l.idx]
}

func (l *Lexer) nextChar() byte {
	return l.Input[l.idx+1]
}
