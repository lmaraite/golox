package expr

import (
	"github.com/lmaraite/golox/token"
)

type Visitor interface {
	VisitBinaryExpr(b Binary) interface{}
	VisitGroupingExpr(g Grouping) interface{}
	VisitLiteralExpr(l Literal) interface{}
	VisitUnaryExpr(u Unary) interface{}
}

type Expr interface {
	Accept(v Visitor) interface{}
}

type Binary struct {
	Left     Expr
	Operator token.Token
	Right    Expr
}

func (binary Binary) Accept(v Visitor) interface{} {
	return v.VisitBinaryExpr(binary)
}

type Grouping struct {
	Expression Expr
}

func (grouping Grouping) Accept(v Visitor) interface{} {
	return v.VisitGroupingExpr(grouping)
}

type Literal struct {
	Value interface{}
}

func (literal Literal) Accept(v Visitor) interface{} {
	return v.VisitLiteralExpr(literal)
}

type Unary struct {
	Operator token.Token
	Right    Expr
}

func (unary Unary) Accept(v Visitor) interface{} {
	return v.VisitUnaryExpr(unary)
}
