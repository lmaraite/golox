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

func (a AstPrinter) VisitBinaryExpr(b expr.Binary) interface{} {
	return a.paranthesize(b.Operator.Lexeme, b.Left, b.Right)
}

func (a AstPrinter) VisitLiteralExpr(l expr.Literal) interface{} {
	if l.Value == nil {
		return "nil"
	}
	return l.Value
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
