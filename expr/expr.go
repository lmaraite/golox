package expr

import (
	"github.com/lmaraite/golox/token"
)

type Visitor interface {
	VisitBinaryExpr(binary Binary) (interface{}, error)
	VisitGroupingExpr(grouping Grouping) (interface{}, error)
	VisitLiteralExpr(literal Literal) (interface{}, error)
	VisitUnaryExpr(unary Unary) (interface{}, error)
	VisitStmtExpr(stmt Stmt) (interface{}, error)
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

type Stmt struct {
	Expression Expr
	Print      Expr
}

func (stmt Stmt) Accept(v Visitor) (interface{}, error) {
	return v.VisitStmtExpr(stmt)
}
