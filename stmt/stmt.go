package stmt

import (
	"github.com/lmaraite/golox/expr"
	"github.com/lmaraite/golox/token"
)

type Visitor interface {
	VisitBlockStmt(Block) error
	VisitExprStmt(Expr) error
	VisitIfStmt(If) error
	VisitPrintStmt(Print) error
	VisitVarStmt(Var) error
}

type Stmt interface {
	Accept(v Visitor) error
}

type Block struct {
	Statements []Stmt
}

func (b Block) Accept(v Visitor) error {
	return v.VisitBlockStmt(b)
}

type Expr struct {
	Expression expr.Expr
}

func (e Expr) Accept(v Visitor) error {
	return v.VisitExprStmt(e)
}

type If struct {
	Condition  expr.Expr
	ThenBranch Stmt
	ElseBranch Stmt
}

func (i If) Accept(v Visitor) error {
	return v.VisitIfStmt(i)
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
