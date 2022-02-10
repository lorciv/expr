package expr

import (
	"errors"
	"strconv"
)

type parser struct {
	tokens chan Token
	cur    Token
}

func (p *parser) parseExpr() (Expr, error) {
	left, err := p.parseTerm()
	if err != nil {
		return nil, err
	}

Loop:
	for p.cur.Type != TokenEOF {
		switch p.cur.Type {
		case TokenPlus:
			p.cur = <-p.tokens
			right, err := p.parseTerm()
			if err != nil {
				return nil, err
			}
			left = add{left, right}
		case TokenMinus:
			p.cur = <-p.tokens
			right, err := p.parseTerm()
			if err != nil {
				return nil, err
			}
			left = sub{left, right}
		default:
			break Loop
		}
	}

	return left, nil
}

func (p *parser) parseTerm() (Expr, error) {
	left, err := p.parseFactor()
	if err != nil {
		return nil, err
	}

Loop:
	for p.cur.Type != TokenEOF {
		switch p.cur.Type {
		case TokenStar:
			p.cur = <-p.tokens
			right, err := p.parseFactor()
			if err != nil {
				return nil, err
			}
			left = mul{left, right}
		case TokenSlash:
			p.cur = <-p.tokens
			right, err := p.parseFactor()
			if err != nil {
				return nil, err
			}
			left = div{left, right}
		default:
			break Loop
		}
	}

	return left, nil
}

func (p *parser) parseFactor() (Expr, error) {
	if p.cur.Type == TokenNumber {
		v, err := strconv.ParseFloat(p.cur.Value, 64)
		if err != nil {
			return nil, err
		}
		p.cur = <-p.tokens
		return literal(v), nil
	}
	if p.cur.Type == TokenLParen {
		p.cur = <-p.tokens
		inner, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if p.cur.Type != TokenRParen {
			return nil, errors.New("missing (")
		}
		p.cur = <-p.tokens
		return inner, nil
	}
	if p.cur.Type == TokenError {
		return nil, errors.New(p.cur.Value)
	}

	return nil, errors.New("unrecognized")
}

func Parse(input string) (Expr, error) {
	p := parser{
		tokens: Lex(input),
	}
	p.cur = <-p.tokens
	return p.parseExpr()
}
