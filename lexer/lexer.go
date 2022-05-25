package lexer

import (
	"errors"
	"fmt"

	"github.com/lmaraite/golox/token"
)

type lexer struct {
	source  string
	tokens  []token.Token
	start   int
	current int
	line    int
}

func Run(source string) error {
	lexer := newLexer(source)
	tokens, err := lexer.scanTokens(source)
	if err != nil {
		return err
	}
	for _, token := range tokens {
		fmt.Println(token.String())
	}
	return nil
}

func newError(line int, message string) error {
	formattedMessage := fmt.Sprintf("[line %d] Error: %s", line, message)
	return errors.New(formattedMessage)
}

func newLexer(source string) *lexer {
	return &lexer{
		source:  source,
		tokens:  make([]token.Token, 0),
		start:   0,
		current: 0,
		line:    1,
	}
}

func (l *lexer) scanTokens(source string) ([]token.Token, error) {
	var err error
	for !l.isAtEnd() {
		l.start = l.current
		err = l.scanToken()
	}
	l.addToken(token.EOF, nil)
	return l.tokens, err
}

func (l *lexer) scanToken() error {
	c := l.advance()
	switch c {
	case '(':
		l.addToken(token.LEFT_PAREN, nil)
	case ')':
		l.addToken(token.RIGHT_PAREN, nil)
	case '{':
		l.addToken(token.LEFT_BRACE, nil)
	case '}':
		l.addToken(token.RIGHT_BRACE, nil)
	case ',':
		l.addToken(token.COMMA, nil)
	case '.':
		l.addToken(token.DOT, nil)
	case '-':
		l.addToken(token.MINUS, nil)
	case '+':
		l.addToken(token.PLUS, nil)
	case ';':
		l.addToken(token.SEMICOLON, nil)
	case '*':
		l.addToken(token.STAR, nil)
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
	case ' ':
		break
	case '\r':
		break
	case '\t':
		break
	case '\n':
		l.line++
	default:
		return newError(l.line, "unexpected character")
	}
	return nil
}

func (l *lexer) advance() uint8 {
	c := l.source[l.current]
	l.current++
	return c
}

func (l *lexer) peek() uint8 {
	if l.isAtEnd() {
		return 0
	}
	return l.source[l.current]
}

func (l *lexer) addToken(tokenType token.TokenType, literal interface{}) {
	lexeme := l.source[l.start:l.current]
	l.tokens = append(l.tokens, *token.NewToken(tokenType, lexeme, literal, l.line))
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
		l.addToken(token.SLASH, nil)
	}
}

func (l *lexer) lexTwoCharToken(tokenType token.TokenType, equalTokenType token.TokenType) {
	if l.match('=') {
		l.addToken(equalTokenType, nil)
	} else {
		l.addToken(tokenType, nil)
	}
}

func (l *lexer) match(expected uint8) bool {
	if l.isAtEnd() {
		return false
	}
	if l.source[l.current] != expected {
		return false
	}
	l.current++
	return true
}

func (l *lexer) isAtEnd() bool {
	return l.current >= len(l.source)
}
