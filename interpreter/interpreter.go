package interpreter

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/lmaraite/golox/environment"
	"github.com/lmaraite/golox/expr"
	"github.com/lmaraite/golox/stmt"
	"github.com/lmaraite/golox/token"
)

func newError(errorToken token.Token, message string) error {
	var formattedMessage string
	if errorToken.TokenType == token.EOF {
		formattedMessage = fmt.Sprintf("[line %d] Runtime error at end: %s", errorToken.Line, message)
	} else {
		formattedMessage = fmt.Sprintf("[line %d] Runtime error at '%s': %s", errorToken.Line, errorToken.Lexeme, message)
	}
	return errors.New(formattedMessage)
}

type Interpreter struct {
	env *environment.Environment
}

func NewInterpreter() *Interpreter {
	return &Interpreter{
		env: environment.NewEmptyEnvironment(),
	}
}

func (i *Interpreter) Interpret(statements []stmt.Stmt) error {
	for _, statement := range statements {
		err := i.execute(statement)
		if err != nil {
			return err
		}
	}
	return nil
}

func (i *Interpreter) execute(statement stmt.Stmt) error {
	return statement.Accept(i)
}

func (i *Interpreter) executeBlock(statements []stmt.Stmt, env *environment.Environment) error {
	previousEnv := i.env
	i.env = env
	for _, statement := range statements {
		err := i.execute(statement)
		if err != nil {
			i.env = previousEnv
			return err
		}
	}
	i.env = previousEnv
	return nil
}

func (i *Interpreter) VisitBlockStmt(statement stmt.Block) error {
	return i.executeBlock(statement.Statements, environment.NewEnvironment(i.env))
}

func (i *Interpreter) VisitExprStmt(statement stmt.Expr) error {
	_, err := i.Evaluate(statement.Expression)
	return err
}

func (i *Interpreter) VisitPrintStmt(statement stmt.Print) error {
	value, err := i.Evaluate(statement.Expression)
	if err != nil {
		return err
	}
	fmt.Println(value)
	return nil
}

func (i *Interpreter) VisitVarStmt(statement stmt.Var) error {
	if statement.Initializer != nil {
		value, err := i.Evaluate(statement.Initializer)
		if err != nil {
			return err
		}
		i.env.Define(statement.Name.Lexeme, value)
	} else {
		i.env.Define(statement.Name.Lexeme, nil)
	}
	return nil
}

func (i *Interpreter) Evaluate(expression expr.Expr) (interface{}, error) {
	return expression.Accept(i)
}

func (i *Interpreter) VisitBinaryExpr(binary expr.Binary) (interface{}, error) {
	left, err := i.Evaluate(binary.Left)
	if err != nil {
		return nil, err
	}
	right, err := i.Evaluate(binary.Right)
	if err != nil {
		return nil, err
	}

	switch binary.Operator.TokenType {
	case token.GREATER:
		err := checkNumberOperands(binary.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) > right.(float64), nil
	case token.GREATER_EQUAL:
		err := checkNumberOperands(binary.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) >= right.(float64), nil
	case token.LESS:
		err := checkNumberOperands(binary.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) < right.(float64), nil
	case token.LESS_EQUAL:
		err := checkNumberOperands(binary.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) <= right.(float64), nil
	case token.MINUS:
		err := checkNumberOperands(binary.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) - right.(float64), nil
	case token.PLUS:
		if reflect.TypeOf(left).Name() == "float64" && reflect.TypeOf(right).Name() == "float64" {
			return left.(float64) + right.(float64), nil
		}
		if reflect.TypeOf(left).Name() == "string" && reflect.TypeOf(right).Name() == "string" {
			return left.(string) + right.(string), nil
		}
		return nil, newError(binary.Operator, "Operands must be two numbers or two strings.")
	case token.SLASH:
		return left.(float64) / right.(float64), nil
	case token.STAR:
		return left.(float64) * right.(float64), nil
	case token.BANG_EQUAL:
		return !isEqual(left, right), nil
	case token.EQUAL_EQUAL:
		return isEqual(left, right), nil
	}
	return nil, nil
}

func (i *Interpreter) VisitGroupingExpr(grouping expr.Grouping) (interface{}, error) {
	return i.Evaluate(grouping.Expression)
}

func (i *Interpreter) VisitLiteralExpr(literal expr.Literal) (interface{}, error) {
	return literal.Value, nil
}

func (i *Interpreter) VisitUnaryExpr(unary expr.Unary) (interface{}, error) {
	right, err := i.Evaluate(unary.Right)
	if err != nil {
		return nil, err
	}

	switch unary.Operator.TokenType {
	case token.BANG:
		return !right.(bool), nil
	case token.MINUS:
		err := checkNumberOperand(unary.Operator, right)
		if err != nil {
			return nil, err
		}
		return -right.(float64), nil
	}
	return nil, nil // unreachable
}

func (i *Interpreter) VisitVariableExpr(variable expr.Variable) (interface{}, error) {
	return i.env.Get(variable.Name)
}

func (i *Interpreter) VisitAssignExpr(assign expr.Assign) (interface{}, error) {
	value, err := i.Evaluate(assign.Value)
	if err != nil {
		return nil, err
	}
	i.env.Assign(assign.Name, value)
	return value, nil
}

func checkNumberOperand(operator token.Token, operand interface{}) error {
	if reflect.TypeOf(operand).String() == "float64" {
		return nil
	}
	return newError(operator, "Operands must be a numbers.")
}

func checkNumberOperands(operator token.Token, left, right interface{}) error {
	if reflect.TypeOf(left).String() == "float64" && reflect.TypeOf(right).String() == "float64" {
		return nil
	}
	return newError(operator, "Operands must be a numbers.")
}

func isEqual(a, b interface{}) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil {
		return false
	}
	return a == b
}
