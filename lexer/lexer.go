package lexer

import (
	"errors"
	"fmt"
	"strconv"
	"unicode"

	"github.com/lmaraite/golox/token"
)

type lexer struct {
	source  string
	tokens  []token.Token
	start   int
	current int
	line    int
}

func NewLexer(source string) *lexer {
	return &lexer{
		source:  source,
		tokens:  make([]token.Token, 0),
		start:   0,
		current: 0,
		line:    1,
	}
}

func newError(line int, message string) error {
	formattedMessage := fmt.Sprintf("[line %d] Error: %s", line, message)
	return errors.New(formattedMessage)
}

func (l *lexer) ScanTokens(source string) ([]token.Token, error) {
	var err error
	for !l.isAtEnd() {
		l.start = l.current
		err = l.scanToken()
	}
	l.addToken(token.EOF)
	return l.tokens, err
}

func (l *lexer) scanToken() error {
	c := l.advance()
	switch c {
	case '(':
		l.addToken(token.LEFT_PAREN)
	case ')':
		l.addToken(token.RIGHT_PAREN)
	case '{':
		l.addToken(token.LEFT_BRACE)
	case '}':
		l.addToken(token.RIGHT_BRACE)
	case ',':
		l.addToken(token.COMMA)
	case '.':
		l.addToken(token.DOT)
	case '-':
		l.addToken(token.MINUS)
	case '+':
		l.addToken(token.PLUS)
	case ';':
		l.addToken(token.SEMICOLON)
	case '*':
		l.addToken(token.STAR)
	case '!':
		l.lexTwoCharToken(token.BANG, token.BANG_EQUAL)
	case '=':
		l.lexTwoCharToken(token.EQUAL, token.EQUAL_EQUAL)
	case '<':
		l.lexTwoCharToken(token.LESS, token.LESS_EQUAL)
	case '>':
		l.lexTwoCharToken(token.GREATER, token.GREATER_EQUAL)
	case '/':
		l.lexSlashOrComment()
	case '"':
		return l.lexString()
	case ' ':
		break
	case '\r':
		break
	case '\t':
		break
	case '\n':
		l.line++
	default:
		if isDigit(c) {
			return l.lexNumber()
		} else if isAlphaNumeric(c) {
			l.lexIdentifier()
		}
		return newError(l.line, "unexpected character")
	}
	return nil
}

// addToken adds a token of a certain tokenType to the lexer
func (l *lexer) addToken(tokenType token.TokenType) {
	lexeme := l.source[l.start:l.current]
	l.tokens = append(l.tokens, *token.NewToken(tokenType, lexeme, nil, l.line))
}

// addLiteralToken adds a literal token to the lexer
func (l *lexer) addLiteralToken(tokenType token.TokenType, literal interface{}) {
	lexeme := l.source[l.start:l.current]
	l.tokens = append(l.tokens, *token.NewToken(tokenType, lexeme, literal, l.line))
}

// advance consumes one character of the lexer's source string and returns it
func (l *lexer) advance() uint8 {
	c := l.source[l.current]
	l.current++
	return c
}

// peek returns the current character of the lexer's source string without consuming it
func (l *lexer) peek() uint8 {
	if l.isAtEnd() {
		return 0
	}
	return l.source[l.current]
}

// peekNext returns the next character of the lexer's source string without consuming any character
func (l *lexer) peekNext() uint8 {
	if l.current+1 >= len(l.source) {
		return 0
	}
	return l.source[l.current+1]
}

func isDigit(char uint8) bool {
	return unicode.IsDigit(rune(char))
}

func isAlpha(c uint8) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}

func isAlphaNumeric(c uint8) bool {
	return isAlpha(c) || isDigit(c)
}

func (l *lexer) lexString() error {
	for l.peek() != '"' && !l.isAtEnd() {
		if l.peek() == '\n' {
			l.line++
		}
		l.advance()
	}
	if l.isAtEnd() {
		return newError(l.line, "unterminated string")
	}
	l.advance() // the closing "

	value := l.source[l.start+1 : l.current-1]
	l.addLiteralToken(token.STRING, value)
	return nil
}

func (l *lexer) lexNumber() error {
	for isDigit(l.peek()) {
		l.advance()
	}
	// Look for a fractional part
	if l.peek() == '.' && isDigit(l.peekNext()) {
		// consume the '.'
		l.advance()
	}
	for isDigit(l.peek()) {
		l.advance()
	}
	value, err := strconv.ParseFloat(l.source[l.start:l.current], 64)
	l.addLiteralToken(token.NUMBER, value)
	return err
}

func (l *lexer) lexIdentifier() {
	for isAlphaNumeric(l.peek()) {
		l.advance()
	}
	text := l.source[l.start:l.current]
	if tokenType, ok := token.Keywords[text]; ok {
		l.addToken(tokenType)
	} else {
		l.addToken(token.IDENTIFIER)
	}

}

func (l *lexer) lexSlashOrComment() {
	if l.match('/') {
		// We have an inline comment,
		// so we need to consume the rest
		// of the line
		for l.peek() != '\n' && !l.isAtEnd() {
			l.advance()
		}
	} else {
		l.addToken(token.SLASH)
	}
}

func (l *lexer) lexTwoCharToken(tokenType token.TokenType, equalTokenType token.TokenType) {
	if l.match('=') {
		l.addToken(equalTokenType)
	} else {
		l.addToken(tokenType)
	}
}

// match checks if the current character of the lexer source string matches the expected character
func (l *lexer) match(expected uint8) bool {
	if l.isAtEnd() {
		return false
	}
	if l.peek() != expected {
		return false
	}
	l.current++
	return true
}

func (l *lexer) isAtEnd() bool {
	return l.current >= len(l.source)
}
