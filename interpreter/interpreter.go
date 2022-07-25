package interpreter

import (
	"reflect"

	"github.com/lmaraite/golox/expr"
	"github.com/lmaraite/golox/token"
)

type Interpreter struct {
}

func (i *Interpreter) evaluate(expression expr.Expr) interface{} {
	return expression.Accept(i)
}

func (i *Interpreter) VisitBinaryExpr(binary expr.Binary) interface{} {
	left := i.evaluate(binary.Left)
	right := i.evaluate(binary.Right)

	switch binary.Operator.TokenType {
	case token.GREATER:
		return left.(float32) > right.(float32)
	case token.GREATER_EQUAL:
		return left.(float32) >= right.(float32)
	case token.LESS:
		return left.(float32) < right.(float32)
	case token.LESS_EQUAL:
		return left.(float32) <= right.(float32)
	case token.MINUS:
		return left.(float32) - right.(float32)
	case token.PLUS:
		if reflect.TypeOf(left).Name() == "float32" && reflect.TypeOf(right).Name() == "float32" {
			return left.(float32) + right.(float32)
		}
		if reflect.TypeOf(left).Name() == "string" && reflect.TypeOf(right).Name() == "string" {
			return left.(string) + right.(string)
		}
	case token.SLASH:
		return left.(float32) / right.(float32)
	case token.STAR:
		return left.(float32) * right.(float32)
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
	return i.evaluate(grouping.Expression)
}

func (i *Interpreter) VisitLiteralExpr(literal expr.Literal) interface{} {
	return literal.Value
}

func (i *Interpreter) VisitUnaryExpr(unary expr.Unary) interface{} {
	right := i.evaluate(unary.Right)

	switch unary.Operator.TokenType {
	case token.BANG:
		return !right.(bool)
	case token.MINUS:
		return -right.(float32)
	}
	return nil // unreachable
}
