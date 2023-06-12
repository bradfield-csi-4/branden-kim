package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

var in_error bool = false

func reportError(line int, message string, code string) {
	fmt.Fprintln(os.Stderr, "ERROR: on line "+strconv.Itoa(line)+" due to "+message+" at code: "+code)
	in_error = true
}

func scanFile(path string) *Scanner {
	dat, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	sc := new(Scanner)
	sc.source = string(dat)
	sc.tokens = make([]*Token, 0)
	sc.start = 0
	sc.current = 0
	sc.line = 1

	return sc
}

func scanInput() *Scanner {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	err := scanner.Err()
	if err != nil {
		fmt.Printf("Error due to: %s\n", err)
		panic(err)
	}

	sc := new(Scanner)
	sc.source = scanner.Text()
	sc.tokens = make([]*Token, 0)
	sc.start = 0
	sc.current = 0
	sc.line = 1

	return sc
}

func main() {
	// if len(os.Args) > 2 {
	// 	fmt.Println("Too many arguments! Please pass in a text file or nothing to paste a query in!")
	// 	os.Exit(1)
	// } else if len(os.Args) == 2 {
	// 	scanner := scanFile(os.Args[1])
	// 	scanner.scanTokens()
	// 	for _, token := range scanner.tokens {
	// 		fmt.Println(token.toString())
	// 	}
	// } else {
	// 	fmt.Print(">> ")
	// 	scanner := scanInput()
	// 	scanner.scanTokens()
	// 	for _, token := range scanner.tokens {
	// 		fmt.Println(token.toString())
	// 	}
	// }

	expression := newBinaryExpr(Token{tokentype: AND, lexeme: "AND"}, newUnaryExpr(Token{tokentype: NOT, lexeme: "NOT"}, newLiteral(Token{tokentype: TERM, lexeme: "123"})), newUnaryExpr(Token{tokentype: NOT, lexeme: "NOT"}, newLiteral(Token{tokentype: TERM, lexeme: "1231231"})))

	fmt.Println(new(ASTVistor).print(expression))
}
