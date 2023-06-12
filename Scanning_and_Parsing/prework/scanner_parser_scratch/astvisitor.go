package main

type Visitor interface {
	visitBinaryExpr(*BinaryExpr) string
	visitUnaryExpr(*UnaryExpr) string
	visitLiteralExpr(*Literal) string
}

type ASTVistor struct {
}

func (v *ASTVistor) print(expr Expr) string {
	return expr.accept(v)
}

func (v *ASTVistor) visitBinaryExpr(expr *BinaryExpr) string {
	return v.parenthesize(expr.operator.lexeme, expr.left, expr.right)
}

func (v *ASTVistor) visitUnaryExpr(e *UnaryExpr) string {
	return v.parenthesize(e.operator.lexeme, e.right)
}

func (v *ASTVistor) visitLiteralExpr(expr *Literal) string {
	if len(expr.value) == 0 {
		return "nil"
	}
	return expr.value
}

func (v *ASTVistor) parenthesize(lexeme string, exprs ...Expr) string {
	var return_string string = ""
	return_string = "(" + lexeme

	for _, expr := range exprs {
		return_string += " " + expr.accept(v)
	}

	return_string += ")"

	return return_string
}
