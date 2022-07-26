package parser

import (
	"errors"
	"fmt"

	"github.com/lmaraite/golox/expr"
	"github.com/lmaraite/golox/token"
)

// This is the context-free grammar we can parse with this parser:
// program        → statement* EOF ;
// statement      → exprStmt
//                | printStmt ;
// exprStmt       → expression ";" ;
// printStmt      → "print" expression ";" ;
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

func (p *parser) Parse() ([]expr.Stmt, error) {
	var statements []expr.Stmt
	for !p.isAtEnd() {
		stmt, err := p.statement()
		if err != nil {
			return nil, err
		}
		statements = append(statements, stmt)
	}
	return statements, nil
}

// statement → exprStmt
//           | printStmt ;
func (p *parser) statement() (expr.Stmt, error) {
	if p.match(token.PRINT) {
		return p.printStatement()
	}
	return p.expressionStatement()
}

// printStmt → "print" expression ";" ;
func (p *parser) printStatement() (expr.Stmt, error) {
	value, err := p.expression()
	if err != nil {
		return expr.Stmt{}, err
	}
	_, err = p.consume(token.SEMICOLON, "Expect ';' after value.")
	return expr.Stmt{Expression: value}, err
}

// exprStmt → expression ";" ;
func (p *parser) expressionStatement() (expr.Stmt, error) {
	expression, err := p.expression()
	if err != nil {
		return expr.Stmt{}, err
	}
	_, err = p.consume(token.SEMICOLON, "Expect ';' after expression.")
	return expr.Stmt{Expression: expression}, err
}

// expression → equality ;
func (p *parser) expression() (expr.Expr, error) {
	return p.equality()
}

// equality → comparison ( ( "!=" | "==" ) comparison )* ;
func (p *parser) equality() (expr.Expr, error) {
	expression, err := p.comparison()
	if err != nil {
		return nil, err
	}
	for p.match(token.BANG_EQUAL, token.EQUAL_EQUAL) {
		operator := p.previous()
		right, err := p.comparison()
		if err != nil {
			return nil, err
		}
		expression = expr.Binary{
			Left:     expression,
			Operator: operator,
			Right:    right,
		}
	}
	return expression, nil
}

// comparison → term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
func (p *parser) comparison() (expr.Expr, error) {
	expression, err := p.term()
	if err != nil {
		return nil, err
	}
	for p.match(token.GREATER, token.GREATER_EQUAL, token.LESS, token.LESS_EQUAL) {
		operator := p.previous()
		right, err := p.term()
		if err != nil {
			return nil, err
		}
		expression = expr.Binary{
			Left:     expression,
			Operator: operator,
			Right:    right,
		}
	}
	return expression, nil
}

// term → factor ( ( "-" | "+" ) factor )* ;
func (p *parser) term() (expr.Expr, error) {
	expression, err := p.factor()
	if err != nil {
		return nil, err
	}
	for p.match(token.MINUS, token.PLUS) {
		operator := p.previous()
		right, err := p.factor()
		if err != nil {
			return nil, err
		}
		expression = expr.Binary{
			Left:     expression,
			Operator: operator,
			Right:    right,
		}
	}
	return expression, nil
}

// factor → unary ( ( "/" | "*" ) unary )* ;
func (p *parser) factor() (expr.Expr, error) {
	expression, err := p.unary()
	if err != nil {
		return nil, err
	}
	for p.match(token.SLASH, token.STAR) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		expression = expr.Binary{
			Left:     expression,
			Operator: operator,
			Right:    right,
		}
	}
	return expression, nil
}

// unary → ( "!" | "-" ) unary
//       | primary ;
func (p *parser) unary() (expr.Expr, error) {
	if p.match(token.BANG, token.MINUS) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		return expr.Unary{
			Operator: operator,
			Right:    right,
		}, nil
	}
	return p.primary()
}

// primary → NUMBER | STRING | "true" | "false" | "nil"
//         | "(" expression ")" ;
func (p *parser) primary() (expr.Expr, error) {
	if p.match(token.FALSE) {
		return expr.Literal{Value: false}, nil
	}
	if p.match(token.TRUE) {
		return expr.Literal{Value: true}, nil
	}
	if p.match(token.NIL) {
		return expr.Literal{Value: nil}, nil
	}
	if p.match(token.NUMBER, token.STRING) {
		return expr.Literal{
			Value: p.previous().Literal,
		}, nil
	}
	if p.match(token.LEFT_PAREN) {
		expression, err := p.expression()
		if err != nil {
			return nil, err
		}
		p.consume(token.RIGHT_PAREN, "Expected ')' after expression.")
		return expr.Grouping{Expression: expression}, nil
	}
	return nil, newError(p.peek(), "Expected expression.")
}

func (p *parser) consume(tokenType token.TokenType, errMsg string) (token.Token, error) {
	if p.check(tokenType) {
		return p.advance(), nil
	}
	return p.advance(), newError(p.peek(), errMsg)
}

func (p *parser) synchronize() {
	p.advance()
	for !p.isAtEnd() {
		if p.previous().TokenType == token.SEMICOLON {
			return
		}
		switch p.peek().TokenType {
		case token.CLASS:
		case token.FUN:
		case token.FOR:
		case token.IF:
		case token.WHILE:
		case token.PRINT:
		case token.RETURN:
			return
		}
		p.advance()
	}
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
