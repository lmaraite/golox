package astprinter

import (
	"fmt"

	"github.com/lmaraite/golox/expr"
)

type AstPrinter struct {
}

func (a AstPrinter) Print(e expr.Expr) (interface{}, error) {
	return e.Accept(a)
}

func (a AstPrinter) VisitBinaryExpr(binary expr.Binary) (interface{}, error) {
	return a.paranthesize(binary.Operator.Lexeme, binary.Left, binary.Right), nil
}

func (a AstPrinter) VisitGroupingExpr(grouping expr.Grouping) (interface{}, error) {
	return a.paranthesize("group", grouping.Expression), nil
}

func (a AstPrinter) VisitLiteralExpr(literal expr.Literal) (interface{}, error) {
	if literal.Value == nil {
		return "nil", nil
	}
	return literal.Value, nil
}

func (a AstPrinter) VisitLogicalExpr(logical expr.Logical) (interface{}, error) {
	return a.paranthesize(logical.Operator.Lexeme, logical.Left, logical.Right), nil
}

func (a AstPrinter) VisitUnaryExpr(unary expr.Unary) (interface{}, error) {
	return a.paranthesize(unary.Operator.Lexeme, unary.Right), nil
}

func (a AstPrinter) VisitVariableExpr(variable expr.Variable) (interface{}, error) {
	return variable.Name.Lexeme, nil
}

func (a AstPrinter) VisitAssignExpr(assign expr.Assign) (interface{}, error) {
	return a.paranthesize(assign.Name.Lexeme, assign.Value), nil
}

func (a AstPrinter) paranthesize(name string, expressions ...expr.Expr) string {
	var result string
	result = fmt.Sprintf("(%s", name)
	for _, v := range expressions {
		print, _ := v.Accept(a)
		result = fmt.Sprintf("%s %v", result, print)
	}
	result = fmt.Sprintf("%s)", result)
	return result
}
