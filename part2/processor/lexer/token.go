package lexer

import "fmt"

type Type string
type Value string

// token type
const (
	LEFT_CURLY_BRACKET   Type = "LEFT_CURLY_BRACKET"
	RIGHT_CURLY_BRACKET  Type = "RIGHT_CURLY_BRACKET"
	LEFT_SQUARE_BRACKET  Type = "LEFT_SQUARE_BRACKET"
	RIGHT_SQUARE_BRACKET Type = "RIGHT_SQUARE_BRACKET"
	COLON                Type = "COLON"
	COMMA                Type = "COMMA"

	NUMBER Type = "NUMBER"
	STRING Type = "STRING"
	TRUE   Type = "TRUE"
	FALSE  Type = "FALSE"
	NULL   Type = "NULL"

	UNKNOWN Type = "UNKNOWN"
)

type Token struct {
	Type  Type
	Value Value
}

type Tokens []Token

func (t Tokens) String() string {
	str := ""
	for _, token := range t {
		str += fmt.Sprintf("(%s, %s)\n", token.Type, token.Value)
	}
	return str
}
