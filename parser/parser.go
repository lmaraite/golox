package parser

import (
	"errors"
	"fmt"

	"github.com/lmaraite/golox/expr"
	"github.com/lmaraite/golox/token"
)

// This is the context-free grammar we can parse with this parser:
// expression     → equality ;
// equality       → comparison ( ( "!=" | "==" ) comparison )* ;
// comparison     → term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
// term           → factor ( ( "-" | "+" ) factor )* ;
// factor         → unary ( ( "/" | "*" ) unary )* ;
// unary          → ( "!" | "-" ) unary
//                | primary ;
// primary        → NUMBER | STRING | "true" | "false" | "nil"
//                | "(" expression ")" ;
type parser struct {
	tokens  []token.Token
	current int
}

func NewParser(tokens []token.Token) *parser {
	return &parser{
		tokens:  tokens,
		current: 0,
	}
}

func newError(errorToken token.Token, message string) error {
	var formattedMessage string
	if errorToken.TokenType == token.EOF {
		formattedMessage = fmt.Sprintf("[line %d] Error at end: %s", errorToken.Line, message)
	} else {
		formattedMessage = fmt.Sprintf("[line %d] Error at '%s': %s", errorToken.Line, errorToken.Lexeme, message)
	}
	return errors.New(formattedMessage)
}

func (p *parser) Parse() expr.Expr {
	return p.expression()
}

// expression → equality ;
func (p *parser) expression() expr.Expr {
	return p.equality()
}

// equality → comparison ( ( "!=" | "==" ) comparison )* ;
func (p *parser) equality() expr.Expr {
	expression := p.comparison()
	for p.match(token.BANG_EQUAL, token.EQUAL_EQUAL) {
		operator := p.previous()
		right := p.comparison()
		expression = expr.Binary{
			Left:     expression,
			Operator: operator,
			Right:    right,
		}
	}
	return expression
}

// comparison → term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
func (p *parser) comparison() expr.Expr {
	expression := p.term()
	for p.match(token.GREATER, token.GREATER_EQUAL, token.LESS, token.LESS_EQUAL) {
		operator := p.previous()
		right := p.term()
		expression = expr.Binary{
			Left:     expression,
			Operator: operator,
			Right:    right,
		}
	}
	return expression
}

// term → factor ( ( "-" | "+" ) factor )* ;
func (p *parser) term() expr.Expr {
	expression := p.factor()
	for p.match(token.MINUS, token.PLUS) {
		operator := p.previous()
		right := p.factor()
		expression = expr.Binary{
			Left:     expression,
			Operator: operator,
			Right:    right,
		}
	}
	return expression
}

// factor → unary ( ( "/" | "*" ) unary )* ;
func (p *parser) factor() expr.Expr {
	expression := p.unary()
	for p.match(token.SLASH, token.STAR) {
		operator := p.previous()
		right := p.unary()
		expression = expr.Binary{
			Left:     expression,
			Operator: operator,
			Right:    right,
		}
	}
	return expression
}

// unary → ( "!" | "-" ) unary
//       | primary ;
func (p *parser) unary() expr.Expr {
	if p.match(token.BANG, token.MINUS) {
		operator := p.previous()
		right := p.unary()
		return expr.Unary{
			Operator: operator,
			Right:    right,
		}
	}
	return p.primary()
}

// primary → NUMBER | STRING | "true" | "false" | "nil"
//         | "(" expression ")" ;
func (p *parser) primary() expr.Expr {
	if p.match(token.FALSE) {
		return expr.Literal{Value: false}
	}
	if p.match(token.TRUE) {
		return expr.Literal{Value: true}
	}
	if p.match(token.NIL) {
		return expr.Literal{Value: nil}
	}
	if p.match(token.NUMBER, token.STRING) {
		return expr.Literal{
			Value: p.previous().Literal,
		}
	}
	if p.match(token.LEFT_PAREN) {
		expression := p.expression()
		p.consume(token.RIGHT_PAREN, "Expected ')' after expression.")
		return expr.Grouping{Expression: expression}
	}
	panic(newError(p.peek(), "Expect expression."))
}

func (p *parser) consume(tokenType token.TokenType, errMsg string) token.Token {
	if p.check(tokenType) {
		return p.advance()
	}
	panic(newError(p.peek(), errMsg))
}

// match if the current token has any of the given types. If so
// it consumes token and returns true.
func (p *parser) match(tokenTypes ...token.TokenType) bool {
	for _, v := range tokenTypes {
		if p.check(v) {
			p.advance()
			return true
		}
	}
	return false
}

// check returns true if the current token is of the given type
func (p *parser) check(tokenType token.TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().TokenType == tokenType
}

// advance consumes the current token and returns it
func (p *parser) advance() token.Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

// isAtEnd checks if weve run out of tokens to parse
func (p *parser) isAtEnd() bool {
	return p.peek().TokenType == token.EOF
}

// peek return the current token without consuming it
func (p *parser) peek() token.Token {
	return p.tokens[p.current]
}

// previous returns the most recently consumed token
func (p *parser) previous() token.Token {
	return p.tokens[p.current-1]
}
