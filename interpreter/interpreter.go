package interpreter

import (
	"reflect"

	"github.com/lmaraite/golox/expr"
	"github.com/lmaraite/golox/token"
)

type Interpreter struct {
}

func (i *Interpreter) Evaluate(expression expr.Expr) interface{} {
	return expression.Accept(i)
}

func (i *Interpreter) VisitBinaryExpr(binary expr.Binary) interface{} {
	left := i.Evaluate(binary.Left)
	right := i.Evaluate(binary.Right)

	switch binary.Operator.TokenType {
	case token.GREATER:
		return left.(float64) > right.(float64)
	case token.GREATER_EQUAL:
		return left.(float64) >= right.(float64)
	case token.LESS:
		return left.(float64) < right.(float64)
	case token.LESS_EQUAL:
		return left.(float64) <= right.(float64)
	case token.MINUS:
		return left.(float64) - right.(float64)
	case token.PLUS:
		if reflect.TypeOf(left).Name() == "float64" && reflect.TypeOf(right).Name() == "float64" {
			return left.(float64) + right.(float64)
		}
		if reflect.TypeOf(left).Name() == "string" && reflect.TypeOf(right).Name() == "string" {
			return left.(string) + right.(string)
		}
	case token.SLASH:
		return left.(float64) / right.(float64)
	case token.STAR:
		return left.(float64) * right.(float64)
	case token.BANG_EQUAL:
		return !isEqual(left, right)
	case token.EQUAL_EQUAL:
		return isEqual(left, right)
	}
	return nil
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

func (i *Interpreter) VisitGroupingExpr(grouping expr.Grouping) interface{} {
	return i.Evaluate(grouping.Expression)
}

func (i *Interpreter) VisitLiteralExpr(literal expr.Literal) interface{} {
	return literal.Value
}

func (i *Interpreter) VisitUnaryExpr(unary expr.Unary) interface{} {
	right := i.Evaluate(unary.Right)

	switch unary.Operator.TokenType {
	case token.BANG:
		return !right.(bool)
	case token.MINUS:
		return -right.(float64)
	}
	return nil // unreachable
}
