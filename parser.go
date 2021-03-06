package expr

import (
	"errors"
	"fmt"
	"strconv"
)

type parser struct {
	tokens    chan Token
	cur, next Token
}

func (p *parser) advance() {
	p.cur = p.next
	p.next = <-p.tokens
}

func (p *parser) parseExpr() (Expr, error) {
	left, err := p.parseTerm()
	if err != nil {
		return nil, err
	}

	for {
		if p.cur.Type == TokenPlus {
			p.advance()
			right, err := p.parseTerm()
			if err != nil {
				return nil, err
			}
			left = add{binary{left, right}}
		} else if p.cur.Type == TokenMinus {
			p.advance()
			right, err := p.parseTerm()
			if err != nil {
				return nil, err
			}
			left = sub{binary{left, right}}
		} else {
			break
		}
	}

	return left, nil
}

func (p *parser) parseTerm() (Expr, error) {
	left, err := p.parseFactor()
	if err != nil {
		return nil, err
	}

	for {
		if p.cur.Type == TokenStar {
			p.advance()
			right, err := p.parseFactor()
			if err != nil {
				return nil, err
			}
			left = mul{binary{left, right}}
		} else if p.cur.Type == TokenSlash {
			p.advance()
			right, err := p.parseFactor()
			if err != nil {
				return nil, err
			}
			left = div{binary{left, right}}
		} else {
			break
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
		p.advance()
		return literal(v), nil
	}
	if p.cur.Type == TokenIdent && p.next.Type == TokenLParen {
		c := call{fn: p.cur.Value}
		p.advance()
		p.advance()

		e, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		c.args = []Expr{e}

		for p.cur.Type == TokenComma {
			p.advance()
			e, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			c.args = append(c.args, e)
		}

		if p.cur.Type != TokenRParen {
			return nil, errors.New("missing )")
		}

		p.advance()
		return c, nil
	}
	if p.cur.Type == TokenIdent {
		v := Var(p.cur.Value)
		p.advance()
		return v, nil
	}
	if p.cur.Type == TokenLParen {
		p.advance()
		inner, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if p.cur.Type != TokenRParen {
			return nil, errors.New("missing (")
		}
		p.advance()
		return inner, nil
	}
	if p.cur.Type == TokenMinus {
		p.advance()
		x, err := p.parseTerm()
		if err != nil {
			return nil, err
		}
		return neg{x: x}, nil
	}
	if p.cur.Type == TokenError {
		return nil, errors.New(p.cur.Value)
	}

	return nil, errors.New("unrecognized")
}

// Parse parses an expression. The parsed expression can then be evaluated.
func Parse(input string) (Expr, error) {
	p := parser{
		tokens: Lex(input),
	}
	p.cur = <-p.tokens
	p.next = <-p.tokens
	e, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if p.cur.Type != TokenEOF {
		if p.cur.Type == TokenError {
			return nil, errors.New(p.cur.Value)
		}
		return nil, errors.New("invalid expression")
	}
	return e, nil
}

// Eval parses, checks and evaluates an expression in the given environment, and returns the result. All variables
// must be explicitly set in the environment.
//
// Eval is a practical alternative to Parse for one-off calculations. If multiple evaluations of the same
// expression are required, consider using Parse once and then calling the Eval method on the parsed expression
// as needed in order to avoid repeating the parsing upon every evaluation.
func Eval(input string, env Env) (res float64, err error) {
	e, err := Parse(input)
	if err != nil {
		return 0, err
	}

	vars := make(map[Var]bool)
	if err := e.Check(vars); err != nil {
		return 0, err
	}
	for v := range vars {
		if _, ok := env[v]; !ok {
			return 0, fmt.Errorf("missing var %s in environment", v)
		}
	}

	return e.Eval(env), nil
}
