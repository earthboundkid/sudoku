//Yet another simple Sudoku solver.
//This one is a translation from Python.
//Original is http://jakevdp.github.io/blog/2013/04/15/code-golf-in-python-sudoku/
package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
)

//The ConnectionGraph is always the same, so it's initialized globally.
var Graph = NewConnectionGraph()

func main() {
	var p Puzzle

	r := bufio.NewReader(os.Stdin)

	for {
		line, err := r.ReadSlice('\n')
		if err != nil {
			break
		}
		err = p.ReadInput(line)

		if err != nil {
			fmt.Println(err.Error())
			break
		}
		//Only gets first solution found.
		p.Solve()
		fmt.Println(&p)

	}

	//	p, _ := NewPuzzle("027800061000030008910005420500016030000970200070000096700000080006027000030480007")
	//	p, _ = NewPuzzle("000000000000003085001020000000507000004000100090000000500000073002010000000040009")
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
//You shouldn't need it unless you improve that.
func (graph *ConnectionGraph) Print() {
	for i := range graph {
		fmt.Printf("%2d: [", i)
		for j := range graph[i] {
			fmt.Printf("%3d", graph[i][j])
		}
		fmt.Print("]\n")
	}

}

//A seen set is a bitfield list of which values have been seen in a space.
type SeenSet uint16

//Rule out one possibility.
func See(s SeenSet, c uint8) SeenSet {
	if c == 0 {
		return s
	}
	return s | 1<<c
}

//How many possibilities are left in the set?
func (s SeenSet) Left() (count int) {
	for i := uint16(1); i < 10; i++ {
		if s&(1<<i) == 0 {
			count += 1
		}
	}
	return count
}

//A puzzle is a 9x9 array
type Puzzle [81]uint8

//Makes a new Puzzle or dies trying!
func (p *Puzzle) ReadInput(input []byte) error {
	if len(input) < 81 {
		return errors.New("Input is too small.")
	}
	for i := range p {
		if input[i] == '.' {
			p[i] = 0
			continue
		}
		if '0' > input[i] || input[i] > '9' {
			return errors.New("Input should only have numbers 0-9.")
		}
		p[i] = input[i] - '0'
	}
	return nil
}

//Solve(c) starts a goroutine that writes solutions to itself 
//to the channel it return and closes channel when done.
func (p *Puzzle) Solve() bool {
	var minPossIndex, minPossCount int = -1, 0
	var minSeen SeenSet

	for i := range p {
		if p[i] == 0 {
			var s SeenSet

			//If it's not filled, go through all the connected squares
			//and eliminate those as possibilities.
			for _, connectionIndex := range Graph[i] {
				s = See(s, p[connectionIndex])
			}

			possCount := s.Left()

			switch {
			//If it wasn't anything, something's wrong with the puzzle, give up.
			case possCount == 0:
				return false
			//If it's the smallest possibilities left we've seen yet, 
			//then save this set for later.
			case minPossCount > possCount:
				fallthrough
			//This is the first zero we've seen, so make it the minimum set.
			case minPossIndex == -1:
				minPossIndex = i
				minPossCount = possCount
				minSeen = s
			}
		}
	}

	//If there were no zeros, then this is a solution, so we're done.
	if minPossIndex == -1 {
		return true
	}

	//OK, let's try out each of the possibilities, and see if any of them
	//solve the problem for us.
	for n := uint8(1); n < 10; n++ {
		if minSeen&(1<<n) != 0 {
			continue
		}

		p[minPossIndex] = n

		if p.Solve() {
			return true
		}
	}

	//We must have barked up the wrong tree. Give up this slot, start over.
	p[minPossIndex] = 0
	return false
}

//Make a new puzzle with just one part different.
func (p Puzzle) Modify(index int, to int) Puzzle {
	p[index] = byte(to) + 1
	return p
}

//Pretty prints a Sudoku puzzle.
func (p *Puzzle) Print() {
	const (
		fRow    = "%s|%s|%s\n"
		divider = "---+---+---\n"
	)
	s := p.String()
	for i := 0; i < 81; i += 27 {
		for j := 0; j < 9; j += 3 {
			fmt.Printf(fRow, s[i+j+0:i+j+3], s[i+j+3:i+j+6], s[i+j+6:i+j+9])
		}
		if i < 54 {
			fmt.Print(divider)
		}
	}
}

//Just dumps it as a single string. Use .Print() for pretty printing.
func (p *Puzzle) String() string {
	b := make([]byte, 81)
	for i := range b {
		b[i] = p[i] + '0'
	}
	return string(b)
}
