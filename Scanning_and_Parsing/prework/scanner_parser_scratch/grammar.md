# CFG for the example language we are implementing

statement  ->   statement
            |   literal
            |   unary
            |   binary

literal    ->   IDENTIFIER | "nil"

unary      ->   ("-" | "NOT") expression

binary     ->   expression operator expression

operator   ->   "AND" | "OR"

NOT should be right associative

# After making the CFG non ambiguous

expression  ->  binary

binary      ->  unary (("AND" | "OR") unary) *

unary       -> ("-" | "NOT") | unary | literal

literal     ->  expression | IDENTIFIER | "nil"