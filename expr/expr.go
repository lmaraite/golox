package expr

import (
	"github.com/lmaraite/golox/token"
)

type Visitor interface {
	VisitBinaryExpr(b Binary) (interface{}, error)
	VisitGroupingExpr(g Grouping) (interface{}, error)
	VisitLiteralExpr(l Literal) (interface{}, error)
	VisitUnaryExpr(u Unary) (interface{}, error)
}

type Expr interface {
	Accept(v Visitor) (interface{}, error)
}

type Binary struct {
	Left     Expr
	Operator token.Token
	Right    Expr
}

func (binary Binary) Accept(v Visitor) (interface{}, error) {
	return v.VisitBinaryExpr(binary)
}

type Grouping struct {
	Expression Expr
}

func (grouping Grouping) Accept(v Visitor) (interface{}, error) {
	return v.VisitGroupingExpr(grouping)
}

type Literal struct {
	Value interface{}
}

func (literal Literal) Accept(v Visitor) (interface{}, error) {
	return v.VisitLiteralExpr(literal)
}

type Unary struct {
	Operator token.Token
	Right    Expr
}

func (unary Unary) Accept(v Visitor) (interface{}, error) {
	return v.VisitUnaryExpr(unary)
}
