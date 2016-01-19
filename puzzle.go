// Package sudoku provides types for solving and displaying sudoku puzzles
package sudoku

import "fmt"

type Digit uint16

//A puzzle is a 9x9 array of bitflag digits.
type Puzzle [81]Digit

func (p *Puzzle) ReadInput(input []byte) error {
	if len(input) < 81 {
		return fmt.Errorf("Input is too small.")
	}
	for i := range p {
		if input[i] == '.' || input[i] == '0' {
			p[i] = 0
			continue
		}
		if '0' > input[i] || input[i] > '9' {
			return fmt.Errorf("Input should only have numbers 0-9.")
		}
		p[i] = 1 << uint16(input[i]-'0')
	}
	return nil
}

func (p *Puzzle) Solve() error {
	if !p.solved() {
		return fmt.Errorf("Couldn't solve the puzzle!")
	}
	return nil
}

func (p *Puzzle) solved() bool {
	var (
		minCount Digit = 0xFFFF
		minFlags Digit
		minIndex int
	)

	for i := range p {
		if p[i] == 0 {
			var f Digit //Flags digits seen.

			//If it's not filled, go through all the connected squares
			//and eliminate those as possibilities.
			for _, connectionIndex := range Graph[i] {
				//Is there anyway to speed this up?
				f |= p[connectionIndex]
			}

			//We eliminated all possibilities. This must be a bad solution try.
			if f == 0x03FE {
				return false
			}

			//Count digits seen
			var c Digit
			f = ^f //Flip digits.
			for i := uint16(1); i < 10; i++ {
				//If the digit wasn't flagged before, it's 1 now.
				//So, scoot it to the 0th place, & out the other bits, and add that.
				c += 1 & (f >> i)
			}

			//Doesn't have fewer constraints than another we saw, try another
			if c >= minCount {
				continue
			}

			//Fewest possibile values to explore, so save it for later
			minCount = c
			minFlags = f
			minIndex = i

			//If it only had one, this is as good as it gets. Move on.
			if c == 1 {
				break
			}
		}
	}

	//If there were no zeros, then this is a solution, so we're done.
	if minCount == 0xFFFF {
		return true
	}

	//OK, let's try out each of the possibilities, and see if any of them
	//solve the problem for us.
	for n := Digit(1); n < 10; n++ {
		if minFlags&(1<<n) != 0 {
			p[minIndex] = 1 << n

			if p.solved() {
				return true
			}
		}

	}

	//We must have barked up the wrong tree. Give up this slot, start over.
	p[minIndex] = 0
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
