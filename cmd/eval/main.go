package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/lorciv/expr"
)

var (
	from = flag.Float64("from", 0, "min value for x")
	to   = flag.Float64("to", 1, "max value for x")
	step = flag.Float64("step", 0.1, "step from min to max")
)

func main() {
	log.SetPrefix("")
	log.SetFlags(0)

	flag.Parse()

	if flag.NArg() < 1 {
		log.Fatalf("Usage: %s expr", os.Args[0])
	}

	e, err := expr.Parse(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}

	vars := make(map[expr.Var]bool)
	if err := e.Check(vars); err != nil {
		log.Fatal(err)
	}
	for v := range vars {
		if v != "x" {
			log.Fatalf("undefined variable: %s", v)
		}
	}

	for x := *from; x < *to; x += *step {
		fmt.Printf("%f\t%f\n", x, e.Eval(expr.Env{"x": x}))
	}
}
