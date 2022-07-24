package expression

import (
	"github.com/lmaraite/golox/token"
)

type Visitor interface {
	visitBinaryExpr(expr Binary) interface{}
}

type Expr interface {
	accept(v Visitor) interface{}
}

type Binary struct {
	left     Expr
	operator token.Token
	right    Expr
}

func (b Binary) accept(v Visitor) interface{} {
	return v.visitBinaryExpr(b)
}
