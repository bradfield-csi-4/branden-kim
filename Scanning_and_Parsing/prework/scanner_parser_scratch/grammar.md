# CFG for the example language we are implementing

statement  ->   statement
            |   literal
            |   unary
            |   binary

literal    ->   IDENTIFIER | "nil"

unary      ->   ("-" | "NOT") expression

binary     ->   expression operator expression

operator   ->   "AND" | "OR"