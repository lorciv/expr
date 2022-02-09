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

type token struct {
	typ   tokenType
	value string
}

func (t token) String() string {
	des := [...]string{"EOF", "LP", "RP", "PLUS", "MINUS", "STAR", "SLASH", "NUM"}
	return fmt.Sprintf("{%s %q}", des[t.typ], t.value)
}

const eof = -1

type lexer struct {
	input  string     // the string being scanned
	start  int        // start position of current token
	pos    int        // current position in the input
	size   int        // size of the last rune read from input
	tokens chan token // channel of scanned tokens (output)
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
	l.tokens <- token{
		typ:   typ,
		value: l.input[l.start:l.pos],
	}
	l.start = l.pos
}

type stateFn func(s *lexer) stateFn

func lexToken(l *lexer) stateFn {
	for {
		r := l.next()
		if r == eof {
			l.emit(tokenEOF)
			return nil
		}
		if !unicode.IsSpace(r) {
			l.backup()
			l.ignore()
			break
		}
	}

	r := l.next()
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

func (l *lexer) run() {
	state := lexToken
	for state != nil {
		state = state(l)
	}
	close(l.tokens)
}

func lex(input string) chan token {
	l := &lexer{
		input:  input,
		tokens: make(chan token),
	}
	go l.run() // Concurrently run the state machine
	return l.tokens
}

func main() {
	scan := bufio.NewScanner(os.Stdin)
	for scan.Scan() {
		for t := range lex(scan.Text()) {
			fmt.Println(t)
		}
	}
	if err := scan.Err(); err != nil {
		log.Fatal(err)
	}
}
