//Yet another simple Sudoku solver.
//This one is a translation from Python.
//Original is http://jakevdp.github.io/blog/2013/04/15/code-golf-in-python-sudoku/
package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/carlmjohnson/sudoku"
)

func ordie(err error) {
	if err != nil {
		fmt.Printf("Failure: %s\n", err)
		os.Exit(-1)
	}
}

func main() {
	var p sudoku.Puzzle

	r := bufio.NewReader(os.Stdin)

	for {
		line, err := r.ReadSlice('\n')
		if err != nil {
			break
		}

		err = p.ReadInput(line)
		ordie(err)
		err = p.Solve()
		ordie(err)
		fmt.Println(&p)
	}
}
