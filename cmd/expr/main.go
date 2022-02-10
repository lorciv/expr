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
		fmt.Println(expr.Parse(scan.Text()))
	}
	if err := scan.Err(); err != nil {
		log.Fatal(err)
	}
}
