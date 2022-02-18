# Expr

Expr is an expression parser and evaluator written in Go.

## References

Here are the main resources that I found helpful.

### Lexical Scanning in Go (video) by Rob Pike

Inspiration for the Lexer. Also, elegant implementation of a state machine.

YouTube: [https://youtu.be/HxaD_trXwRE](https://youtu.be/HxaD_trXwRE)

### Recursive descent parsing (video) by Harry Porter

Simple and clear explanation of the recursive descent parsing algorithm, from theory to implementation.

YouTube: [https://youtu.be/SToUyjAsaFk](https://youtu.be/SToUyjAsaFk)

### The Go Programming Language (book) by Alan Donovan, Brian Kernighan

Chapter 7.9 contains an example of how to represent the abstract syntax tree of a parsed expression,
as well as how to evaulate its result.

Website: [http://www.gopl.io](http://www.gopl.io)

### Crafting interpreters (book) by Robert Nystrom

A guide to implement a programming language interpreter from start to end, including lexing, parsing and evaluation.

Website: [http://craftinginterpreters.com](http://craftinginterpreters.com)
