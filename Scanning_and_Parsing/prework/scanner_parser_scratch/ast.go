package main

type Expr interface {
	accept(Visitor) string
}

type UnaryExpr struct {
	right    Expr
	operator Token
}

func newUnaryExpr(operator Token, expr Expr) *UnaryExpr {
	new_expr := &UnaryExpr{expr, operator}
	return new_expr
}

func (e *UnaryExpr) accept(v Visitor) string {
	return v.visitUnaryExpr(e)
}

type Literal struct {
	value string
}

func newLiteral(operator Token) *Literal {
	e := &Literal{operator.lexeme}
	return e
}

func (e *Literal) accept(v Visitor) string {
	return v.visitLiteralExpr(e)
}

type BinaryExpr struct {
	left     Expr
	right    Expr
	operator Token
}

func newBinaryExpr(operator Token, left_expr Expr, right_expr Expr) *BinaryExpr {
	new_expr := &BinaryExpr{left_expr, right_expr, operator}
	return new_expr
}

func (e *BinaryExpr) accept(v Visitor) string {
	return v.visitBinaryExpr(e)
}
