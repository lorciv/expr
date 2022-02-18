package expr

import (
	"fmt"
	"math"
)

// Env maps variables to their respective numeric value. It represents an environment in which expressions can be evaluated.
type Env map[Var]float64

// Expr is a parsed arithmetic expression that can be checked and evaluated. To avoid the risk of run-time errors,
// checking should be done before evaulation.
//
// Eval evaluates the expression in the given environment. The environment can be used to set the value of the variables
// that appear in the expression. Unset variables default to 0.
//
// Check checks that the expression is valid and returns the set of defined variables. A valid expression
// can be safely evaluated without run-time errors. The initial empty set of variables must be provided by the caller,
// due to the recursive nature of this function.
type Expr interface {
	Eval(env Env) float64
	Check(vars map[Var]bool) error
}

// Var is an expression consisting of a named variable. It is exported so that clients can populate the environment (Env).
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

type binary struct {
	left, right Expr
}

func (b binary) Check(vars map[Var]bool) error {
	if err := b.left.Check(vars); err != nil {
		return err
	}
	return b.right.Check(vars)
}

type add struct {
	binary
}

func (a add) Eval(env Env) float64 {
	return a.left.Eval(env) + a.right.Eval(env)
}

func (a add) String() string {
	return fmt.Sprintf("{add %v %v}", a.left, a.right)
}

type sub struct {
	binary
}

func (s sub) Eval(env Env) float64 {
	return s.left.Eval(env) - s.right.Eval(env)
}

func (s sub) String() string {
	return fmt.Sprintf("{sub %v %v}", s.left, s.right)
}

type mul struct {
	binary
}

func (m mul) Eval(env Env) float64 {
	return m.left.Eval(env) * m.right.Eval(env)
}

func (m mul) String() string {
	return fmt.Sprintf("{mul %v %v}", m.left, m.right)
}

type div struct {
	binary
}

func (d div) Eval(env Env) float64 {
	return d.left.Eval(env) / d.right.Eval(env)
}

func (d div) String() string {
	return fmt.Sprintf("{div %v %v}", d.left, d.right)
}

type call struct {
	fn   string
	args []Expr
}

func (c call) Eval(env Env) float64 {
	switch c.fn {
	case "sin":
		return math.Sin(c.args[0].Eval(env))
	case "cos":
		return math.Cos(c.args[0].Eval(env))
	case "pow":
		return math.Pow(c.args[0].Eval(env), c.args[1].Eval(env))
	}
	panic("unknown function call: " + c.fn)
}

var arity = map[string]int{"sin": 1, "cos": 1, "pow": 2}

func (c call) Check(vars map[Var]bool) error {
	a, ok := arity[c.fn]
	if !ok {
		return fmt.Errorf("unknown function call: %s", c.fn)
	}
	if len(c.args) != a {
		return fmt.Errorf("call to %s has %d args, expected %d", c.fn, len(c.args), a)
	}
	for _, e := range c.args {
		if err := e.Check(vars); err != nil {
			return err
		}
	}
	return nil
}

func (c call) String() string {
	return fmt.Sprintf("{call %s %v}", c.fn, c.args)
}
