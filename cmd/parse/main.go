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
		fmt.Println("Parsed:", e)

		vars := make(map[expr.Var]bool)
		if err := e.Check(vars); err != nil {
			log.Print(err)
			continue
		}
		fmt.Println("Vars:", vars, "(all set to 0)")

		fmt.Println("Eval:", e.Eval(make(expr.Env)))
	}
	if err := scan.Err(); err != nil {
		log.Fatal(err)
	}
}
