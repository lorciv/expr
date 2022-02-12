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
			fmt.Fprintln(os.Stderr, err)
		}
		fmt.Println(e)
	}
	if err := scan.Err(); err != nil {
		log.Fatal(err)
	}
}
