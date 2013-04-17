//Yet another simple Sudoku solver.
//This one is a translation from Python.
//Original is http://jakevdp.github.io/blog/2013/04/15/code-golf-in-python-sudoku/
package main

import (
	"errors"
	"fmt"
)

//The ConnectionGraph is always the same, so it's initialized globally.
var Graph = NewConnectionGraph()

func main() {
	p, _ := NewPuzzle("027800061000030008910005420500016030000970200070000096700000080006027000030480007")
	c := make(chan *Puzzle)
	go p.Solve(c)
	for solution := range c {
		solution.Print()
		fmt.Print("\n~~~ ~~~ ~~~\n\n")
	}
	p, _ = NewPuzzle("000000000000000000000000000000000000000000000000000000000000000000000000000000000")
	c = make(chan *Puzzle)
	go p.Solve(c)
	for solution := range c {
		solution.Print()
		fmt.Print("\n~~~ ~~~ ~~~\n\n")
	}
}

//Each square in a puzzle is connected to 20 other squares, not counting itself
type ConnectionGraph [81][20]int

func NewConnectionGraph() (graph ConnectionGraph) {
	//We're going to go through each square in the graph (i)
	//and every other square (s); if the pair (i and s) are connected,
	//we add it to the graph's inner array then advance the inner array's index
	for i := range graph {
		var j int
		for s := 0; s < 81; s++ {
			//Don't add yourself
			if i == s {
				continue
			}
			switch icol, scol := i%9, s%9; {
			//Same column
			case icol == scol:
				fallthrough
			//Same row
			case i-icol <= s && s < i-icol+9:
				fallthrough
			//Same box
			case i/27 == s/27 && icol/3 == scol/3:
				graph[i][j] = s
				j++
			}
		}
	}
	return graph
}

//Pretty prints what's connected to what.
//This was only needed for debugging the new connection graph alogrithm.
func (graph *ConnectionGraph) Print() {
	for i := range graph {
		fmt.Printf("%2d: [", i)
		for j := range graph[i] {
			fmt.Printf("%3d", graph[i][j])
		}
		fmt.Print("]\n")
	}

}

//A possibility set is a list of which values are still possible for a space.
type PossibilitySet [9]bool

//Sets all fields true
func NewPossibilitySet() *PossibilitySet {
	r := new(PossibilitySet)
	for i := range r {
		r[i] = true
	}
	return r
}

//Rule out one possibility.
func (p *PossibilitySet) Eliminate(c byte) {
	if c == '0' {
		return
	}
	p[c-'1'] = false
}

//How many possibilities are left in the set?
func (p *PossibilitySet) Count() (count int) {
	for _, v := range p {
		if v {
			count += 1
		}
	}
	return count
}

//A puzzle is a 9x9 array
type Puzzle [81]byte

//Makes a new Puzzle or dies trying!
func NewPuzzle(s string) (Puzzle, error) {
	var p Puzzle

	if len(s) != 81 {
		return p, errors.New("Input is the wrong size.")
	}
	for i := range s {
		if '0' > s[i] || s[i] > '9' {
			return p, errors.New("Input should only have numbers 0-9.")
		}
		p[i] = s[i]
	}
	return p, nil
}

//Solve(c) writes solutions to itself to the channel. Closes channel when done.
func (p *Puzzle) Solve(c chan *Puzzle) {
	//We need to close the channel when we're done, so listeners know to 
	//stop listening to us. Downside: can't reuse the channel.
	defer close(c)

	var minPossIndex, minPossCount int
	var minPoss *PossibilitySet

	for i := range p {
		if p[i] == '0' {
			poss := NewPossibilitySet()

			//If it's not filled, go through all the connected squares
			//and eliminate those as possibilities.
			for _, connectionIndex := range Graph[i] {
				poss.Eliminate(p[connectionIndex])
			}

			possCount := poss.Count()

			switch {
			//If it wasn't anything, something's wrong, give up.
			case possCount == 0:
				return
			//If it's the smallest set we've seen yet, 
			//then save this possibility for later.
			case minPossCount > possCount:
				fallthrough
			//This is the first zero we've seen, so make it the minimum set.
			case minPoss == nil:
				minPossIndex = i
				minPossCount = possCount
				minPoss = poss
			}
		}
	}

	//If there were no zeros, then this is a solution, so we're done.
	if minPoss == nil {
		c <- p
		return
	}

	//OK, let's try out each of the possibilities, and see if any of them
	//solve the problem for us.
	for pos := range minPoss {
		if minPoss[pos] == false {
			continue
		}
		np := p.Modify(minPossIndex, pos)
		nc := make(chan *Puzzle)
		go np.Solve(nc)
		for solution := range nc {
			c <- solution
		}
	}
}

//Make a new puzzle with just one part different.
func (p Puzzle) Modify(index int, to int) Puzzle {
	p[index] = byte(to) + '1'
	return p
}

//Pretty prints a Sudoku puzzle.
func (p *Puzzle) Print() {
	const (
		fRow    = "%s|%s|%s\n"
		divider = "---+---+---\n"
	)
	for i := 0; i < 81; i += 27 {
		for j := 0; j < 9; j += 3 {
			fmt.Printf(fRow, p[i+j+0:i+j+3], p[i+j+3:i+j+6], p[i+j+6:i+j+9])
		}
		if i < 54 {
			fmt.Print(divider)
		}
	}
}

//Just dumps it as a single string. Use .Print() for pretty printing.
func (p *Puzzle) String() string {
	b := make([]byte, 0, len(p))
	for _, c := range p {
		b = append(b, c)
	}
	return string(b)
}
