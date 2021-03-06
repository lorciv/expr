package expr

import (
	"fmt"
	"unicode"
	"unicode/utf8"
)

// TokenType is the type of a token.
type TokenType int

// Valid values for TokenType.
const (
	TokenEOF TokenType = iota
	TokenError
	TokenLParen
	TokenRParen
	TokenPlus
	TokenMinus
	TokenStar
	TokenSlash
	TokenNumber
	TokenIdent
	TokenComma
)

// Token is an atom of the expression grammar. It is the product of lexical analysis as returned by Lex.
type Token struct {
	Type  TokenType
	Value string
}

func (t Token) String() string {
	des := [...]string{2: "LPAR", "RPAR", "PLUS", "MINUS", "STAR", "SLASH", "NUM", "IDENT", "COMMA"}
	switch t.Type {
	case TokenEOF:
		return "{EOF}"
	case TokenError:
		return fmt.Sprintf("{ERR %s}", t.Value)
	}
	return fmt.Sprintf("{%s %q}", des[t.Type], t.Value)
}

const eof = -1

type lexer struct {
	input  string     // the string being scanned
	start  int        // start position of current token
	pos    int        // current position in the input
	size   int        // size of the last rune read from input
	tokens chan Token // channel of scanned tokens (output)
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

func (l *lexer) emit(typ TokenType) {
	l.tokens <- Token{
		Type:  typ,
		Value: l.input[l.start:l.pos],
	}
	l.start = l.pos
}

func (l *lexer) errorf(format string, args ...interface{}) stateFn {
	l.tokens <- Token{
		Type:  TokenError,
		Value: fmt.Sprintf(format, args...),
	}
	return nil
}

type stateFn func(s *lexer) stateFn

func lexToken(l *lexer) stateFn {
	for {
		r := l.next()
		if r == eof {
			l.ignore()
			l.emit(TokenEOF)
			return nil
		}
		if !unicode.IsSpace(r) {
			l.backup()
			l.ignore()
			break
		}
	}

	r := l.next()
	if unicode.IsLetter(r) {
		return lexIdentifier
	}
	if unicode.IsDigit(r) {
		return lexNumber
	}

	switch r {
	case '(':
		l.emit(TokenLParen)
		return lexToken
	case ')':
		l.emit(TokenRParen)
		return lexToken
	case '+':
		l.emit(TokenPlus)
		return lexToken
	case '-':
		l.emit(TokenMinus)
		return lexToken
	case '*':
		l.emit(TokenStar)
		return lexToken
	case '/':
		l.emit(TokenSlash)
		return lexToken
	case ',':
		l.emit(TokenComma)
		return lexToken
	}

	return l.errorf("invalid token: %q", l.input[l.start:l.pos])
}

func lexIdentifier(l *lexer) stateFn {
	for unicode.IsLetter(l.peek()) {
		l.next()
	}
	l.emit(TokenIdent)
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
	l.emit(TokenNumber)
	return lexToken
}

func (l *lexer) run() {
	state := lexToken
	for state != nil {
		state = state(l)
	}
	close(l.tokens)
}

// Lex launches lexical analysis on the given expression. It returns a channel of tokens from which
// the caller can retrieve the tokens. When there are no more tokens, the channel will produce tokens of type TokenEOF.
// In case of error, the channel produces a TokenError, with Value containing the error message.
func Lex(input string) chan Token {
	l := &lexer{
		input:  input,
		tokens: make(chan Token),
	}
	go l.run() // Concurrently run the state machine
	return l.tokens
}
