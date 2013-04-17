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

//The ConnectionGraph is always the same, so it's initialized just once globally.
var Graph = ConnectionGraph{{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 18, 19, 20, 27, 36, 45, 54, 63, 72}, {0, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 18, 19, 20, 28, 37, 46, 55, 64, 73}, {0, 1, 3, 4, 5, 6, 7, 8, 9, 10, 11, 18, 19, 20, 29, 38, 47, 56, 65, 74}, {0, 1, 2, 4, 5, 6, 7, 8, 12, 13, 14, 21, 22, 23, 30, 39, 48, 57, 66, 75}, {0, 1, 2, 3, 5, 6, 7, 8, 12, 13, 14, 21, 22, 23, 31, 40, 49, 58, 67, 76}, {0, 1, 2, 3, 4, 6, 7, 8, 12, 13, 14, 21, 22, 23, 32, 41, 50, 59, 68, 77}, {0, 1, 2, 3, 4, 5, 7, 8, 15, 16, 17, 24, 25, 26, 33, 42, 51, 60, 69, 78}, {0, 1, 2, 3, 4, 5, 6, 8, 15, 16, 17, 24, 25, 26, 34, 43, 52, 61, 70, 79}, {0, 1, 2, 3, 4, 5, 6, 7, 15, 16, 17, 24, 25, 26, 35, 44, 53, 62, 71, 80}, {0, 1, 2, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 27, 36, 45, 54, 63, 72}, {0, 1, 2, 9, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 28, 37, 46, 55, 64, 73}, {0, 1, 2, 9, 10, 12, 13, 14, 15, 16, 17, 18, 19, 20, 29, 38, 47, 56, 65, 74}, {3, 4, 5, 9, 10, 11, 13, 14, 15, 16, 17, 21, 22, 23, 30, 39, 48, 57, 66, 75}, {3, 4, 5, 9, 10, 11, 12, 14, 15, 16, 17, 21, 22, 23, 31, 40, 49, 58, 67, 76}, {3, 4, 5, 9, 10, 11, 12, 13, 15, 16, 17, 21, 22, 23, 32, 41, 50, 59, 68, 77}, {6, 7, 8, 9, 10, 11, 12, 13, 14, 16, 17, 24, 25, 26, 33, 42, 51, 60, 69, 78}, {6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 17, 24, 25, 26, 34, 43, 52, 61, 70, 79}, {6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 24, 25, 26, 35, 44, 53, 62, 71, 80}, {0, 1, 2, 9, 10, 11, 19, 20, 21, 22, 23, 24, 25, 26, 27, 36, 45, 54, 63, 72}, {0, 1, 2, 9, 10, 11, 18, 20, 21, 22, 23, 24, 25, 26, 28, 37, 46, 55, 64, 73}, {0, 1, 2, 9, 10, 11, 18, 19, 21, 22, 23, 24, 25, 26, 29, 38, 47, 56, 65, 74}, {3, 4, 5, 12, 13, 14, 18, 19, 20, 22, 23, 24, 25, 26, 30, 39, 48, 57, 66, 75}, {3, 4, 5, 12, 13, 14, 18, 19, 20, 21, 23, 24, 25, 26, 31, 40, 49, 58, 67, 76}, {3, 4, 5, 12, 13, 14, 18, 19, 20, 21, 22, 24, 25, 26, 32, 41, 50, 59, 68, 77}, {6, 7, 8, 15, 16, 17, 18, 19, 20, 21, 22, 23, 25, 26, 33, 42, 51, 60, 69, 78}, {6, 7, 8, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 26, 34, 43, 52, 61, 70, 79}, {6, 7, 8, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 35, 44, 53, 62, 71, 80}, {0, 9, 18, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 45, 46, 47, 54, 63, 72}, {1, 10, 19, 27, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 45, 46, 47, 55, 64, 73}, {2, 11, 20, 27, 28, 30, 31, 32, 33, 34, 35, 36, 37, 38, 45, 46, 47, 56, 65, 74}, {3, 12, 21, 27, 28, 29, 31, 32, 33, 34, 35, 39, 40, 41, 48, 49, 50, 57, 66, 75}, {4, 13, 22, 27, 28, 29, 30, 32, 33, 34, 35, 39, 40, 41, 48, 49, 50, 58, 67, 76}, {5, 14, 23, 27, 28, 29, 30, 31, 33, 34, 35, 39, 40, 41, 48, 49, 50, 59, 68, 77}, {6, 15, 24, 27, 28, 29, 30, 31, 32, 34, 35, 42, 43, 44, 51, 52, 53, 60, 69, 78}, {7, 16, 25, 27, 28, 29, 30, 31, 32, 33, 35, 42, 43, 44, 51, 52, 53, 61, 70, 79}, {8, 17, 26, 27, 28, 29, 30, 31, 32, 33, 34, 42, 43, 44, 51, 52, 53, 62, 71, 80}, {0, 9, 18, 27, 28, 29, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 54, 63, 72}, {1, 10, 19, 27, 28, 29, 36, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 55, 64, 73}, {2, 11, 20, 27, 28, 29, 36, 37, 39, 40, 41, 42, 43, 44, 45, 46, 47, 56, 65, 74}, {3, 12, 21, 30, 31, 32, 36, 37, 38, 40, 41, 42, 43, 44, 48, 49, 50, 57, 66, 75}, {4, 13, 22, 30, 31, 32, 36, 37, 38, 39, 41, 42, 43, 44, 48, 49, 50, 58, 67, 76}, {5, 14, 23, 30, 31, 32, 36, 37, 38, 39, 40, 42, 43, 44, 48, 49, 50, 59, 68, 77}, {6, 15, 24, 33, 34, 35, 36, 37, 38, 39, 40, 41, 43, 44, 51, 52, 53, 60, 69, 78}, {7, 16, 25, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 44, 51, 52, 53, 61, 70, 79}, {8, 17, 26, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 51, 52, 53, 62, 71, 80}, {0, 9, 18, 27, 28, 29, 36, 37, 38, 46, 47, 48, 49, 50, 51, 52, 53, 54, 63, 72}, {1, 10, 19, 27, 28, 29, 36, 37, 38, 45, 47, 48, 49, 50, 51, 52, 53, 55, 64, 73}, {2, 11, 20, 27, 28, 29, 36, 37, 38, 45, 46, 48, 49, 50, 51, 52, 53, 56, 65, 74}, {3, 12, 21, 30, 31, 32, 39, 40, 41, 45, 46, 47, 49, 50, 51, 52, 53, 57, 66, 75}, {4, 13, 22, 30, 31, 32, 39, 40, 41, 45, 46, 47, 48, 50, 51, 52, 53, 58, 67, 76}, {5, 14, 23, 30, 31, 32, 39, 40, 41, 45, 46, 47, 48, 49, 51, 52, 53, 59, 68, 77}, {6, 15, 24, 33, 34, 35, 42, 43, 44, 45, 46, 47, 48, 49, 50, 52, 53, 60, 69, 78}, {7, 16, 25, 33, 34, 35, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 53, 61, 70, 79}, {8, 17, 26, 33, 34, 35, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 62, 71, 80}, {0, 9, 18, 27, 36, 45, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 72, 73, 74}, {1, 10, 19, 28, 37, 46, 54, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 72, 73, 74}, {2, 11, 20, 29, 38, 47, 54, 55, 57, 58, 59, 60, 61, 62, 63, 64, 65, 72, 73, 74}, {3, 12, 21, 30, 39, 48, 54, 55, 56, 58, 59, 60, 61, 62, 66, 67, 68, 75, 76, 77}, {4, 13, 22, 31, 40, 49, 54, 55, 56, 57, 59, 60, 61, 62, 66, 67, 68, 75, 76, 77}, {5, 14, 23, 32, 41, 50, 54, 55, 56, 57, 58, 60, 61, 62, 66, 67, 68, 75, 76, 77}, {6, 15, 24, 33, 42, 51, 54, 55, 56, 57, 58, 59, 61, 62, 69, 70, 71, 78, 79, 80}, {7, 16, 25, 34, 43, 52, 54, 55, 56, 57, 58, 59, 60, 62, 69, 70, 71, 78, 79, 80}, {8, 17, 26, 35, 44, 53, 54, 55, 56, 57, 58, 59, 60, 61, 69, 70, 71, 78, 79, 80}, {0, 9, 18, 27, 36, 45, 54, 55, 56, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74}, {1, 10, 19, 28, 37, 46, 54, 55, 56, 63, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74}, {2, 11, 20, 29, 38, 47, 54, 55, 56, 63, 64, 66, 67, 68, 69, 70, 71, 72, 73, 74}, {3, 12, 21, 30, 39, 48, 57, 58, 59, 63, 64, 65, 67, 68, 69, 70, 71, 75, 76, 77}, {4, 13, 22, 31, 40, 49, 57, 58, 59, 63, 64, 65, 66, 68, 69, 70, 71, 75, 76, 77}, {5, 14, 23, 32, 41, 50, 57, 58, 59, 63, 64, 65, 66, 67, 69, 70, 71, 75, 76, 77}, {6, 15, 24, 33, 42, 51, 60, 61, 62, 63, 64, 65, 66, 67, 68, 70, 71, 78, 79, 80}, {7, 16, 25, 34, 43, 52, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 71, 78, 79, 80}, {8, 17, 26, 35, 44, 53, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 78, 79, 80}, {0, 9, 18, 27, 36, 45, 54, 55, 56, 63, 64, 65, 73, 74, 75, 76, 77, 78, 79, 80}, {1, 10, 19, 28, 37, 46, 54, 55, 56, 63, 64, 65, 72, 74, 75, 76, 77, 78, 79, 80}, {2, 11, 20, 29, 38, 47, 54, 55, 56, 63, 64, 65, 72, 73, 75, 76, 77, 78, 79, 80}, {3, 12, 21, 30, 39, 48, 57, 58, 59, 66, 67, 68, 72, 73, 74, 76, 77, 78, 79, 80}, {4, 13, 22, 31, 40, 49, 57, 58, 59, 66, 67, 68, 72, 73, 74, 75, 77, 78, 79, 80}, {5, 14, 23, 32, 41, 50, 57, 58, 59, 66, 67, 68, 72, 73, 74, 75, 76, 78, 79, 80}, {6, 15, 24, 33, 42, 51, 60, 61, 62, 69, 70, 71, 72, 73, 74, 75, 76, 77, 79, 80}, {7, 16, 25, 34, 43, 52, 60, 61, 62, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 80}, {8, 17, 26, 35, 44, 53, 60, 61, 62, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79}}

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
type Digit uint16

//How many possibilities are left in the set?
func Unflagged(d Digit) (count int) {
	for i := uint16(1); i < 10; i++ {
		if d&(1<<i) == 0 {
			count += 1
		}
	}
	return count
}

//A puzzle is a 9x9 array of bitflag digits.
type Puzzle [81]Digit

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
		p[i] = 1 << uint16(input[i]-'0')
	}
	return nil
}

//Solve(c) starts a goroutine that writes solutions to itself 
//to the channel it return and closes channel when done.
func (p *Puzzle) Solve() bool {
	var minPossIndex, minPossCount int = -1, 0
	var minPossFlags Digit

loop:
	for i := range p {
		if p[i] == 0 {
			var d Digit //Flags digits seen.

			//If it's not filled, go through all the connected squares
			//and eliminate those as possibilities.
			for _, connectionIndex := range Graph[i] {
				//Is there anyway to speed this up?
				d |= p[connectionIndex]
				/* No help with performance: 
				//Is d full?
					if d == 0x3fe {
						break
					}*/
			}

			possCount := Unflagged(d)
			switch {
			//If it wasn't anything, something's wrong with the puzzle, give up.
			case possCount == 0:
				return false
			//Doesn't get more minimaler.
			case possCount == 1:
				minPossIndex = i
				minPossCount = possCount
				minPossFlags = d
				break loop
			//If it's the smallest possibilities left we've seen yet, 
			//then save this set for later.
			case minPossCount > possCount:
				fallthrough
			//This is the first zero we've seen, so make it the minimum set.
			case minPossIndex == -1:
				minPossIndex = i
				minPossFlags = d
				minPossCount = possCount
			}
		}
	}

	//If there were no zeros, then this is a solution, so we're done.
	if minPossIndex == -1 {
		return true
	}

	//OK, let's try out each of the possibilities, and see if any of them
	//solve the problem for us.
	for n := Digit(1); n < 10; n++ {
		if minPossFlags&(1<<n) != 0 {
			continue
		}

		p[minPossIndex] = 1 << n

		if p.Solve() {
			return true
		}
	}

	//We must have barked up the wrong tree. Give up this slot, start over.
	p[minPossIndex] = 0
	return false
}

//Just dumps it as a single string. Use .Print() for pretty printing.
func (p *Puzzle) String() string {
	b := make([]byte, 81)
	for i, v := range p {
		switch v {
		case 0:
			b[i] = '0'
		case 1 << 0:
			b[i] = '0'
		case 1 << 1:
			b[i] = '1'
		case 1 << 2:
			b[i] = '2'
		case 1 << 3:
			b[i] = '3'
		case 1 << 4:
			b[i] = '4'
		case 1 << 5:
			b[i] = '5'
		case 1 << 6:
			b[i] = '6'
		case 1 << 7:
			b[i] = '7'
		case 1 << 8:
			b[i] = '8'
		case 1 << 9:
			b[i] = '9'
		default:
			b[i] = '?'
		}
	}
	return string(b)
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
