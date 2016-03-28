// Package sudoku provides types for solving and displaying sudoku puzzles
package sudoku

import (
	"bytes"
	"fmt"
)

// Digit is a bitflag for representing sudoku digits. 0x0 is unset,
// 0x2 is a one, 0x4 is a two, etc.
type Digit uint16

// Constraints counts the number of flags set on a digit.
func (d Digit) Constraints() (c int) {
	// Fast path for common cases
	if d == 0 {
		return 0
	}

	if d == 0x3fe {
		return 9
	}

	for i := Digit(1); i < 10; i++ {
		// If the constraint was added before, it's 1 now.
		// So, scoot it to the 0th place, & out the other bits, and add that.
		c += int(1 & (d >> i))
	}

	return c
}

// Puzzle is a 9x9 array of bitflag digits.
type Puzzle [81]Digit

// ReadInput sets the puzzle based on byte slice of input.
// Input is expected to be 81 bytes long with 0 or . for empty spaces.
// Input beyond 81 bytes is ignored.
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
		p[i] = 1 << Digit(input[i]-'0')
	}
	return nil
}

// IsValid checks for basic validity (no repeated numbers)
func (p *Puzzle) IsValid() bool {
	for i := range p {
		// Fast path common case
		if p[i] == 0 {
			continue
		}

		for _, j := range Graph[i] {
			if p[i] == p[j] {
				return false
			}
		}
	}
	return true
}

// Solve will mutate the puzzle into a solved state or return an error
// explaining why this was impossible.
func (p *Puzzle) Solve() error {
	if !p.IsValid() {
		return fmt.Errorf("puzzle is not valid")
	}
	if !p.solved() {
		return fmt.Errorf("puzzle is unsolvable")
	}
	return nil
}

func (p *Puzzle) solved() bool {
	var (
		maxConstraints, minIndex = -1, -1
		seen, possibleSolutions  Digit
	)

	for i := range p {
		// Already constrained...
		if p[i] != 0 {
			continue
		}
		// Reset flags
		seen = 0

		// Go through all the connected squares and eliminate those as possibilities.
		for _, connectionIndex := range Graph[i] {
			seen |= p[connectionIndex]
		}

		//Count digits seen
		c := seen.Constraints()

		// We eliminated all possibilities. This must be a bad solution try.
		if c == 9 {
			return false
		}

		// Doesn't have more constraints than another we saw, try another
		if c <= maxConstraints {
			continue
		}

		//Fewest possibile values to explore, so save it for later
		maxConstraints = c
		possibleSolutions = ^seen
		minIndex = i

		//If it only had one possibility left, this is as good as it gets. Move on.
		if c == 8 {
			break
		}
	}

	//If there were no zeros, then this is a solution, so we're done.
	if maxConstraints == -1 {
		return true
	}

	//OK, let's try out each of the possibilities, and see if any of them
	//solve the problem for us.
	for n := Digit(1); n < 10; n++ {
		if v := Digit(1 << n); possibleSolutions&v != 0 {
			p[minIndex] = v

			if p.solved() {
				return true
			}
		}

	}

	//We must have barked up the wrong tree. Give up this slot, start over.
	p[minIndex] = 0
	return false
}

// String dumps a puzzle as a single line. Use .Print() for pretty printing.
func (p *Puzzle) String() string {
	b := make([]byte, 81)
loop:
	for i, v := range p {
		if v == 0 {
			b[i] = '.'
			continue
		}
		for j := uint(0); j < 10; j++ {
			if v == 1<<j {
				b[i] = '0' + byte(j)
				continue loop
			}
		}
		b[i] = '?'
	}
	return string(b)
}

// Print pretty prints a puzzle with dividers etc.
func (p *Puzzle) Print() string {
	const (
		fRow    = "%s|%s|%s\n"
		divider = "---+---+---\n"
	)

	var buf bytes.Buffer
	s := p.String()

	for i := 0; i < 81; i += 9 {
		fmt.Fprintf(&buf, fRow, s[i+0:i+3], s[i+3:i+6], s[i+6:i+9])
		if i%27 == 18 && i < 54 {
			fmt.Fprint(&buf, divider)
		}
	}
	return buf.String()
}
