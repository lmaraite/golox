package parser

import (
	"errors"
	"fmt"

	"github.com/lmaraite/golox/expr"
	"github.com/lmaraite/golox/stmt"
	"github.com/lmaraite/golox/token"
)

// This is the context-free grammar we can parse with this parser:
// program        → declaration* EOF ;
// declaration    → varDecl
//                | statement ;
// varDecl        → "var" IDENTIFIER ( "=" expression )? ";" ;
// statement      → exprStmt
//			      | ifStmt
//                | printStmt
//				  | block ;
// block		  → "{" declaration* "}" ;
// exprStmt       → expression ";" ;
// ifStmt         → "if" "(" expression ")" statement
//                ( "else" statement )? ;
// printStmt      → "print" expression ";" ;
// expression     → assignment ;
// assignment     → IDENTIFIER "=" assignment
//                | equality ;
// equality       → comparison ( ( "!=" | "==" ) comparison )* ;
// comparison     → term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
// term           → factor ( ( "-" | "+" ) factor )* ;
// factor         → unary ( ( "/" | "*" ) unary )* ;
// unary          → ( "!" | "-" ) unary
//                | primary ;
// primary        → "true" | "false" | "nil"
//                | NUMBER | STRING
//                | "(" expression ")"
//                | IDENTIFIER ;
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

func (p *parser) Parse() ([]stmt.Stmt, error) {
	var statements []stmt.Stmt
	for !p.isAtEnd() {
		stmt, err := p.declaration()
		if err != nil {
			return nil, err
		}
		statements = append(statements, stmt)
	}
	return statements, nil
}

// declaration → varDecl
//             | statement ;
func (p *parser) declaration() (stmt.Stmt, error) {
	if p.match(token.VAR) {
		return p.varDeclaration()
	}
	return p.statement()
}

// varDecl → "var" IDENTIFIER ( "=" expression )? ";" ;
func (p *parser) varDeclaration() (stmt.Stmt, error) {
	name, err := p.consume(token.IDENTIFIER, "Expected variable name.")
	if err != nil {
		return nil, err
	}
	var initializer expr.Expr
	if p.match(token.EQUAL) {
		initializer, err = p.expression()
		if err != nil {
			return nil, err
		}
	}
	_, err = p.consume(token.SEMICOLON, "Expected ';' after variable declaration.")
	if err != nil {
		return nil, err
	}
	return stmt.Var{Name: name, Initializer: initializer}, nil
}

// statement → exprStmt
//			 | ifStmt
//           | printStmt
//           | block ;
func (p *parser) statement() (stmt.Stmt, error) {
	if p.match(token.IF) {
		return p.ifStatement()
	}
	if p.match(token.PRINT) {
		return p.printStatement()
	}
	if p.match(token.LEFT_BRACE) {
		statements, err := p.block()
		if err != nil {
			return nil, err
		}
		return stmt.Block{Statements: statements}, nil
	}
	return p.expressionStatement()
}

// ifStmt → "if" "(" expression ")" statement
//          ( "else" statement )? ;
func (p *parser) ifStatement() (stmt.Stmt, error) {
	_, err := p.consume(token.LEFT_PAREN, "Expect '(' after 'if'.")
	if err != nil {
		return nil, err
	}
	condition, err := p.expression()
	if err != nil {
		return nil, err
	}
	_, err = p.consume(token.RIGHT_PAREN, "Expect ')' after if condition.")
	if err != nil {
		return nil, err
	}
	thenBranch, err := p.statement()
	if err != nil {
		return nil, err
	}
	if p.match(token.ELSE) {
		elseBranch, err := p.statement()
		if err != nil {
			return nil, err
		}
		return stmt.If{
			Condition:  condition,
			ThenBranch: thenBranch,
			ElseBranch: elseBranch,
		}, nil
	}
	return stmt.If{
		Condition:  condition,
		ThenBranch: thenBranch,
		ElseBranch: nil,
	}, nil
}

// block → "{" declaration* "}" ;
func (p *parser) block() ([]stmt.Stmt, error) {
	var statements []stmt.Stmt

	for !p.check(token.RIGHT_BRACE) && !p.isAtEnd() {
		statement, err := p.declaration()
		if err != nil {
			return nil, err
		}
		statements = append(statements, statement)
	}

	p.consume(token.RIGHT_BRACE, "Expected '}' after block.")
	return statements, nil
}

// printStmt → "print" expression ";" ;
func (p *parser) printStatement() (stmt.Print, error) {
	value, err := p.expression()
	if err != nil {
		return stmt.Print{}, err
	}
	_, err = p.consume(token.SEMICOLON, "Expect ';' after value.")
	return stmt.Print{Expression: value}, err
}

// exprStmt → expression ";" ;
func (p *parser) expressionStatement() (stmt.Expr, error) {
	expression, err := p.expression()
	if err != nil {
		return stmt.Expr{}, err
	}
	_, err = p.consume(token.SEMICOLON, "Expect ';' after expression.")
	return stmt.Expr{Expression: expression}, err
}

// expression → assignment ;
func (p *parser) expression() (expr.Expr, error) {
	return p.assignment()
}

// assignment → IDENTIFIER "=" assignment
//            | equality ;
func (p *parser) assignment() (expr.Expr, error) {
	expression, err := p.equality()
	if err != nil {
		return nil, err
	}
	if p.match(token.EQUAL) {
		equals := p.previous()
		value, err := p.assignment()
		if err != nil {
			return nil, err
		}
		if variable, ok := expression.(expr.Variable); ok {
			name := variable.Name
			return expr.Assign{Name: name, Value: value}, nil
		}
		return nil, newError(equals, "Invalid assignment target.")
	}
	return expression, nil
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

// primary → "true" | "false" | "nil"
//         | NUMBER | STRING
//         | "(" expression ")"
//         | IDENTIFIER ;
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
	if p.match(token.IDENTIFIER) {
		return expr.Variable{
			Name: p.previous(),
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

// consume consumes a token, and returns it if it matches the tokenType.
// If not, an error is returned.
func (p *parser) consume(tokenType token.TokenType, errMsg string) (token.Token, error) {
	if p.check(tokenType) {
		return p.advance(), nil
	}
	return p.advance(), newError(p.peek(), errMsg)
}

// match if the current token has any of the given types. If so
// it consumes the token and returns true.
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
