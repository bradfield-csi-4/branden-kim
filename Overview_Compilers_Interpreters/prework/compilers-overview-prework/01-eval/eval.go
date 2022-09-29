package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"log"
	"strconv"
)

type ASTVisitor struct {
	stack []string
	index int
}

func (v *ASTVisitor) Visit(node ast.Node) ast.Visitor {
	switch x := node.(type) {
	case *ast.BinaryExpr:
		v.stack[v.index] = x.Op.String()
		v.index += 1
	case *ast.BasicLit:
		v.stack[v.index] = x.Value
		v.index += 1
	}
	return v
}

type Stack []int

// IsEmpty: check if stack is empty
func (s *Stack) IsEmpty() bool {
	return len(*s) == 0
}

// Push a new value onto the stack
func (s *Stack) Push(i int) {
	*s = append(*s, i) // Simply append the new value to the end of the stack
}

// Remove and return top element of stack. Return false if stack is empty.
func (s *Stack) Pop() (int, bool) {
	if s.IsEmpty() {
		return -1, false
	} else {
		index := len(*s) - 1   // Get the index of the top most element.
		element := (*s)[index] // Index into the slice and obtain the element.
		*s = (*s)[:index]      // Remove it from the stack by slicing it off.
		return element, true
	}
}

// Given an expression containing only int types, evaluate
// the expression and return the result.
func Evaluate(expr ast.Expr) (int, error) {
	v := ASTVisitor{make([]string, 100), 0}
	ast.Walk(&v, expr)
	var stack Stack

	// going to have to reverse the contents stored in the stack in order to do a postfix calculator solution
	for i, j := 0, v.index-1; i < j; i, j = i+1, j-1 {
		v.stack[i], v.stack[j] = v.stack[j], v.stack[i]
	}

	for x := 0; x < v.index; x++ {
		token := v.stack[x]

		if val, err := strconv.Atoi(token); err == nil {
			stack.Push(val)
		} else {
			arg1, _ := stack.Pop()
			arg2, _ := stack.Pop()
			if token == "+" {
				stack.Push(arg1 + arg2)
			} else if token == "-" {
				stack.Push(arg1 - arg2)
			} else if token == "*" {
				stack.Push(arg1 * arg2)
			} else if token == "/" {
				stack.Push(arg1 / arg2)
			}
		}
	}
	res, _ := stack.Pop()
	return res, nil
}

func main() {
	expr, err := parser.ParseExpr("1 + 2 - 3 * 4")
	if err != nil {
		log.Fatal(err)
	}
	// fset := token.NewFileSet()
	// err = ast.Print(fset, expr)
	res, err := Evaluate(expr)
	fmt.Println(res)
	if err != nil {
		log.Fatal(err)
	}
}
