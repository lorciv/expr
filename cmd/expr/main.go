package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/lorciv/expr"
)

func main() {
	scan := bufio.NewScanner(os.Stdin)
	for scan.Scan() {
		e, err := expr.Parse(scan.Text())
		if err != nil {
			log.Print(err)
			continue
		}

		vars := make(map[expr.Var]bool)
		if err := e.Check(vars); err != nil {
			log.Print(err)
			continue
		}
		for v := range vars {
			if v != "x" {
				log.Printf("undefined variable %v will be set to 0", v)
			}
		}

		env := make(expr.Env)
		for i := 0; i < 10; i++ {
			env[expr.Var("x")] = float64(i)
			fmt.Println(i, "->", e.Eval(env))
		}
	}
	if err := scan.Err(); err != nil {
		log.Fatal(err)
	}
}
