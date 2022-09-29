package main

import "strconv"

type TokenType int

const (
	// Keywords
	AND TokenType = iota
	OR
	NOT

	// Literals
	TERM

	// EOF
	EOF
)

func (t TokenType) toString() string {
	switch t {
	case AND:
		return "AND"
	case OR:
		return "OR"
	case NOT:
		return "NOT"
	case TERM:
		return "TERM"
	case EOF:
		return "EOF"
	}

	return "unknown"
}

type Token struct {
	tokentype TokenType
	lexeme    string
	line      int
}

func (t Token) toString() string {
	return "Line: " + strconv.Itoa(t.line) + " " + t.tokentype.toString() + " " + t.lexeme
}
