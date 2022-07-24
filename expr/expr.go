package expr

import (
	"github.com/lmaraite/golox/token"
)

type Visitor interface {
	VisitBinaryExpr(b Binary) interface{}
	VisitLiteralExpr(l Literal) interface{}
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

type Literal struct {
	Value interface{}
}

func (l Literal) Accept(v Visitor) interface{} {
	return v.VisitLiteralExpr(l)
}
