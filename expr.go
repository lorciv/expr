package expr

import (
	"fmt"
)

type Env map[Var]float64

// An Expr is an arithmetic expression.
type Expr interface {
	Eval(env Env) float64
	Check(vars map[Var]bool) error
}

type Var string

func (v Var) Eval(env Env) float64 {
	return env[v]
}

func (v Var) Check(vars map[Var]bool) error {
	vars[v] = true
	return nil
}

func (v Var) String() string {
	return fmt.Sprintf("{var %s}", string(v))
}

type literal float64

func (l literal) Eval(_ Env) float64 {
	return float64(l)
}

func (l literal) Check(vars map[Var]bool) error {
	return nil
}

type neg struct {
	x Expr
}

func (n neg) Eval(env Env) float64 {
	return -1 * n.x.Eval(env)
}

func (n neg) Check(vars map[Var]bool) error {
	return n.x.Check(vars)
}

func (n neg) String() string {
	return fmt.Sprintf("{neg %v}", n.x)
}

type add struct {
	left, right Expr
}

func (a add) Eval(env Env) float64 {
	return a.left.Eval(env) + a.right.Eval(env)
}

func (a add) Check(vars map[Var]bool) error {
	if err := a.left.Check(vars); err != nil {
		return err
	}
	return a.right.Check(vars)
}

func (a add) String() string {
	return fmt.Sprintf("{add %v %v}", a.left, a.right)
}

type sub struct {
	left, right Expr
}

func (s sub) Eval(env Env) float64 {
	return s.left.Eval(env) - s.right.Eval(env)
}

func (s sub) Check(vars map[Var]bool) error {
	if err := s.left.Check(vars); err != nil {
		return err
	}
	return s.right.Check(vars)
}

func (s sub) String() string {
	return fmt.Sprintf("{sub %v %v}", s.left, s.right)
}

type mul struct {
	left, right Expr
}

func (m mul) Eval(env Env) float64 {
	return m.left.Eval(env) * m.right.Eval(env)
}

func (m mul) Check(vars map[Var]bool) error {
	if err := m.left.Check(vars); err != nil {
		return err
	}
	return m.right.Check(vars)
}

func (m mul) String() string {
	return fmt.Sprintf("{mul %v %v}", m.left, m.right)
}

type div struct {
	left, right Expr
}

func (d div) Eval(env Env) float64 {
	return d.left.Eval(env) / d.right.Eval(env)
}

func (d div) Check(vars map[Var]bool) error {
	if err := d.left.Check(vars); err != nil {
		return err
	}
	return d.right.Check(vars)
}

func (d div) String() string {
	return fmt.Sprintf("{div %v %v}", d.left, d.right)
}
