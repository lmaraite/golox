package stmt

import "github.com/lmaraite/golox/expr"

type Visitor interface {
	VisitExprStmt(Expr) error
	VisitPrintStmt(Print) error
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
