package stmt

import (
	"github.com/lmaraite/golox/expr"
	"github.com/lmaraite/golox/token"
)

type Visitor interface {
	VisitExprStmt(Expr) error
	VisitPrintStmt(Print) error
	VisitVarStmt(Var) error
}

type Stmt interface {
	Accept(v Visitor) error
}

type Expr struct {
	Expression expr.Expr
}

func (e Expr) Accept(v Visitor) error {
	return v.VisitExprStmt(e)
}

type Print struct {
	Expression expr.Expr
}

func (p Print) Accept(v Visitor) error {
	return v.VisitPrintStmt(p)
}

type Var struct {
	Name        token.Token
	Initializer expr.Expr
}

func (v Var) Accept(vis Visitor) error {
	return vis.VisitVarStmt(v)
}
