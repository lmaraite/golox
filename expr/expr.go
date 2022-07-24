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

func (b Binary) Accept(v Visitor) interface{} {
	return v.VisitBinaryExpr(b)
}

type Grouping struct {
	Expression Expr
}

func (g Grouping) Accept(v Visitor) interface{} {
	return v.VisitGroupingExpr(g)
}

type Literal struct {
	Value interface{}
}

func (l Literal) Accept(v Visitor) interface{} {
	return v.VisitLiteralExpr(l)
}

type Unary struct {
	Operator token.Token
	Right    Expr
}

func (u Unary) Accept(v Visitor) interface{} {
	return v.VisitUnaryExpr(u)
}
