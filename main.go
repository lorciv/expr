package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"unicode"
	"unicode/utf8"
)

type tokenType int

const (
	tokenEOF tokenType = iota

	tokenLParen
	tokenRParen
	tokenPlus
	tokenMinus
	tokenStar
	tokenSlash

	tokenNumber
)

func (t tokenType) String() string {
	switch t {
	case tokenEOF:
		return "eof"
	case tokenLParen:
		return "lparen"
	case tokenRParen:
		return "rparen"
	case tokenPlus:
		return "plus"
	case tokenMinus:
		return "minus"
	case tokenStar:
		return "star"
	case tokenSlash:
		return "slash"
	case tokenNumber:
		return "num"
	}
	return "unknown"
}

type token struct {
	typ   tokenType
	value string
}

func (t token) String() string {
	return fmt.Sprintf("{%v %q}", t.typ, t.value)
}

const eof = -1

type lexer struct {
	input  string  // the string being scanned
	start  int     // start position of current token
	pos    int     // current position in the input
	size   int     // size of the last rune read from input
	tokens []token // list of scanned tokens
}

func (l *lexer) next() rune {
	if l.pos >= len(l.input) {
		l.size = 0
		return eof
	}
	var r rune
	r, l.size = utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += l.size
	return r
}

func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

func (l *lexer) backup() {
	l.pos -= l.size
}

func (l *lexer) ignore() {
	l.start = l.pos
}

func (l *lexer) emit(typ tokenType) {
	t := token{
		typ:   typ,
		value: l.input[l.start:l.pos],
	}
	l.tokens = append(l.tokens, t)
	l.start = l.pos
}

func (s *lexer) run() {
	state := lexToken
	for state != nil {
		state = state(s)
	}
}

type stateFn func(s *lexer) stateFn

func lexToken(l *lexer) stateFn {
	r := l.next()
	if r == eof {
		l.emit(tokenEOF)
		return nil
	}

	if unicode.IsDigit(r) {
		return lexNumber
	}

	switch r {
	case '(':
		l.emit(tokenLParen)
	case ')':
		l.emit(tokenRParen)
	case '+':
		l.emit(tokenPlus)
	case '-':
		l.emit(tokenMinus)
	case '*':
		l.emit(tokenStar)
	case '/':
		l.emit(tokenSlash)
	}
	return lexToken
}

func lexNumber(l *lexer) stateFn {
	for unicode.IsDigit(l.peek()) {
		l.next()
	}
	if l.peek() == '.' {
		l.next()
	}
	for unicode.IsDigit(l.peek()) {
		l.next()
	}
	l.emit(tokenNumber)
	return lexToken
}

func main() {
	scan := bufio.NewScanner(os.Stdin)
	for scan.Scan() {
		l := lexer{
			input: scan.Text(),
		}
		l.run()

		fmt.Println(l.tokens)
	}
	if err := scan.Err(); err != nil {
		log.Fatal(err)
	}
}
