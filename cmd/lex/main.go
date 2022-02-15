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
		for t := range expr.Lex(scan.Text()) {
			fmt.Println(t)
		}
	}
	if err := scan.Err(); err != nil {
		log.Fatal(err)
	}
}
