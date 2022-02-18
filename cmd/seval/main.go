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
		res, err := expr.Eval(scan.Text(), nil)
		if err != nil {
			log.Print(err)
			continue
		}
		fmt.Println(res)
	}
	if err := scan.Err(); err != nil {
		log.Fatal(err)
	}
}
