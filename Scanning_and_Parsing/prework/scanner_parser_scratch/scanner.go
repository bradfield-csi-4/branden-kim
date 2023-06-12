package main

type Scanner struct {
	source  string
	tokens  []*Token
	start   int
	current int
	line    int
}

func (s *Scanner) scanTokens() {
	// Main loop to run through every single character
	for !s.isAtEnd() {
		s.start = s.current

		// Extracting the current byte (character) in the source string
		var current_char string = s.source[s.current : s.current+1]
		s.current += 1

		// Checking all the possible lexemes that the character can consume
		switch current_char {
		case "-":
			s.addToken(NOT)
			break
		// lexemes that don't mean anything and we should skip
		case " ":
			break
		case "\r":
			break
		case "\t":
			break
		case "\n":
			s.line++
			break
		default:
			if s.isAlpha(current_char) {
				s.parseTermOrReserved()
			} else {
				reportError(s.line, "parsing error due to unidentified characters", current_char)
			}
			break
		}
	}

	s.start = s.current
	s.addToken(EOF)
}

func (s *Scanner) addToken(t TokenType) {
	lexeme := s.source[s.start:s.current]

	new_token := new(Token)
	new_token.tokentype = t
	new_token.lexeme = lexeme
	new_token.line = s.line

	new_tokens := append(s.tokens, new_token)
	s.tokens = new_tokens
}

func (s *Scanner) parseTermOrReserved() {
	for s.isAlpha(s.lookaheadOne()) {
		s.current++
	}

	text := s.source[s.start:s.current]
	switch text {
	case "OR":
		s.addToken(OR)
		break
	case "AND":
		s.addToken(AND)
		break
	default:
		s.addToken(TERM)
		break
	}
}

func (s *Scanner) lookaheadOne() string {
	if s.isAtEnd() {
		return "\\0"
	}

	return s.source[s.current : s.current+1]
}

func (s *Scanner) isAlpha(current_char string) bool {
	var char byte = current_char[0]
	return (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char == '_')
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}
