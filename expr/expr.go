package expr

import (
	"github.com/lmaraite/golox/token"
)

type Visitor interface {
	VisitAssignExpr(assign Assign) (interface{}, error)
	VisitBinaryExpr(binary Binary) (interface{}, error)
	VisitGroupingExpr(grouping Grouping) (interface{}, error)
	VisitLiteralExpr(literal Literal) (interface{}, error)
	VisitLogicalExpr(logical Logical) (interface{}, error)
	VisitUnaryExpr(unary Unary) (interface{}, error)
	VisitVariableExpr(variable Variable) (interface{}, error)
}

type Expr interface {
	Accept(visitor Visitor) (interface{}, error)
}

func (a Assign) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitAssignExpr(a)
}

type Assign struct {
	Name  token.Token
	Value Expr
}

type Binary struct {
	Left     Expr
	Operator token.Token
	Right    Expr
}

func (b Binary) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitBinaryExpr(b)
}

type Grouping struct {
	Expression Expr
}

func (g Grouping) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitGroupingExpr(g)
}

type Literal struct {
	Value interface{}
}

func (l Literal) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitLiteralExpr(l)
}

type Logical struct {
	Left     Expr
	Operator token.Token
	Right    Expr
}

func (l Logical) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitLogicalExpr(l)
}

type Unary struct {
	Operator token.Token
	Right    Expr
}

func (u Unary) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitUnaryExpr(u)
}

type Variable struct {
	Name token.Token
}

func (v Variable) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitVariableExpr(v)
}
