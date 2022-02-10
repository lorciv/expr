package expr

import "fmt"

// An Expr is an arithmetic expression.
type Expr interface {
	Eval() float64
	Children() []Expr
}

type literal float64

func (l literal) Eval() float64 {
	return float64(l)
}

func (l literal) Children() []Expr {
	return nil
}

type add struct {
	left, right Expr
}

func (a add) Eval() float64 {
	return a.left.Eval() + a.right.Eval()
}

func (a add) Children() []Expr {
	return []Expr{a.left, a.right}
}

func (a add) String() string {
	return fmt.Sprintf("{add %v %v}", a.left, a.right)
}

type sub struct {
	left, right Expr
}

func (s sub) Eval() float64 {
	return s.left.Eval() - s.right.Eval()
}

func (s sub) Children() []Expr {
	return []Expr{s.left, s.right}
}

func (s sub) String() string {
	return fmt.Sprintf("{sub %v %v}", s.left, s.right)
}

type mul struct {
	left, right Expr
}

func (m mul) Eval() float64 {
	return m.left.Eval() * m.right.Eval()
}

func (m mul) Children() []Expr {
	return []Expr{m.left, m.right}
}

func (m mul) String() string {
	return fmt.Sprintf("{mul %v %v}", m.left, m.right)
}

type div struct {
	left, right Expr
}

func (d div) Eval() float64 {
	return d.left.Eval() / d.right.Eval()
}

func (d div) Children() []Expr {
	return []Expr{d.left, d.right}
}

func (d div) String() string {
	return fmt.Sprintf("{div %v %v}", d.left, d.right)
}
