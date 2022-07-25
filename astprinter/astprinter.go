package astprinter

import (
	"fmt"

	"github.com/lmaraite/golox/expr"
)

type AstPrinter struct {
}

func (a AstPrinter) Print(e expr.Expr) interface{} {
	return e.Accept(a)
}

func (a AstPrinter) VisitBinaryExpr(binary expr.Binary) interface{} {
	return a.paranthesize(binary.Operator.Lexeme, binary.Left, binary.Right)
}

func (a AstPrinter) VisitGroupingExpr(grouping expr.Grouping) interface{} {
	return a.paranthesize("group", grouping.Expression)
}

func (a AstPrinter) VisitLiteralExpr(literal expr.Literal) interface{} {
	if literal.Value == nil {
		return "nil"
	}
	return literal.Value
}

func (a AstPrinter) VisitUnaryExpr(unary expr.Unary) interface{} {
	return a.paranthesize(unary.Operator.Lexeme, unary.Right)
}

func (a AstPrinter) paranthesize(name string, expressions ...expr.Expr) string {
	var result string
	result = fmt.Sprintf("(%s", name)
	for _, v := range expressions {
		result = fmt.Sprintf("%s %v", result, v.Accept(a))
	}
	result = fmt.Sprintf("%s)", result)
	return result
}
