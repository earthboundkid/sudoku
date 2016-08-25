package sudoku_test

import (
	"testing"

	"github.com/carlmjohnson/sudoku"
)

var testcases = []struct {
	name     string
	in       string
	readable bool
	solvable bool
}{
	{"empty string", "", false, false},
	{"empty puzzle", ".................................................................................", true, true},    // Empty case
	{"worst case", "..............3.85..1.2.......5.7.....4...1...9.......5......73..2.1........4...9", true, true},      // near worst case for brute-force solver (wiki)
	{"platinum blonde", ".......12........3..23..4....18....5.6..7.8.......9.....85.....9...4.5..47...6...", true, true}, // gsf's sudoku q1 (Platinum Blonde)
	{"Cheese", ".2..5.7..4..1....68....3...2....8..3.4..2.5.....6...1...2.9.....9......57.4...9..", true, true},          // (Cheese)
	{"fata morgana", "........3..1..56...9..4..7......9.5.7.......8.5.4.2....8..2..9...35..1..6........", true, true},    // (Fata Morgana)
	{"red dwarf", "12.3....435....1....4........54..2..6...7.........8.9...31..5.......9.7.....6...8", true, true},       // (Red Dwarf)
	{"easter monster", "1.......2.9.4...5...6...7...5.9.3.......7.......85..4.7.....6...3...9.8...2.....1", true, true},  // (Easter Monster)
	{"explainer", ".......39.....1..5..3.5.8....8.9...6.7...2...1..4.......9.8..5..2....6..4..7.....", true, true},       // Nicolas Juillerat's Sudoku explainer 1.2.1 (top 5)
	{"", "12.3.....4.....3....3.5......42..5......8...9.6...5.7...15..2......9..6......7..8", true, true},
	{"", "..3..6.8....1..2......7...4..9..8.6..3..4...1.7.2.....3....5.....5...6..98.....5.", true, true},
	{"", "1.......9..67...2..8....4......75.3...5..2....6.3......9....8..6...4...1..25...6.", true, true},
	{"", "..9...4...7.3...2.8...6...71..8....6....1..7.....56...3....5..1.4.....9...2...7..", true, true},
	{"suexrat9", "....9..5..1.....3...23..7....45...7.8.....2.......64...9..1.....8..6......54....7", true, true}, // dukuso's suexrat9 (top 1)
	{"topn87", "4...3.......6..8..........1....5..9..8....6...7.2........1.27..5.3....4.9........", true, true},   // from http://magictour.free.fr/topn87 (top 3)
	{"", "7.8...3.....2.1...5.........4.....263...8.......1...9..9.6....4....7.5...........", true, true},
	{"", "3.7.4...........918........4.....7.....16.......25..........38..9....5...2.6.....", true, true},
	{"suexratt", "........8..3...4...9..2..6.....79.......612...6.5.2.7...8...5...1.....2.4.5.....3", true, true}, // dukuso's suexratt (top 1)
	{"sudoku17", ".......1.4.........2...........5.4.7..8...3....1.9....3..4..2...5.1........8.6...", true, true}, // first 2 from sudoku17
	{"", ".......12....35......6...7.7.....3.....4..8..1...........12.....8.....4..5....6..", true, true},
	{"setbb", "1.......2.9.4...5...6...7...5.3.4.......6........58.4...2...6...3...9.8.7.......1", true, true}, // 2 from http://www.setbb.com/phpbb/viewtopic.php?p=10478
	{"", ".....1.2.3...4.5.....6....7..2.....1.8..9..3.4.....8..5....2....9..3.4....67.....", true, true},
	{"unsolvable", "..1.........345789...789345......................................................", true, false}, // Unsolvable
	{"invalid", "11...............................................................................", true, false},    // Invalid
}

func TestSolver(t *testing.T) {

	for _, test := range testcases {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			var (
				p   sudoku.Puzzle
				err error
			)
			err = p.ReadInput([]byte(test.in))
			if (err == nil) != test.readable {
				t.Errorf("Readability error.\n\tExpected readability = %v\n\tInput = %q.",
					test.readable, test.in)
				return
			}
			if !test.readable {
				return
			}
			// Copy puzzle
			q := p
			solvingErr := p.Solve()
			if test.solvable != (solvingErr == nil) {
				t.Errorf("Solvability error.\n\tExpected solvable = %v\n\tInput = %q.",
					test.solvable, test.in)
				return
			}
			if !test.solvable {
				return
			}
			if !p.IsValid() {
				t.Errorf("Solution invalid.\n\tInput =  %q\n\tOutput = %q.", &q, &p)
			}
			// Check digits match input
			for i := range q {
				if q[i] != 0 && q[i] != p[i] {
					t.Errorf("Solution changed constraints.\nInput =  %q\n\tOutput = %q.", &q, &p)
				}
			}

			t.Logf("\nInput:\n%s\n", q.Print())
			t.Logf("\nOutput:\n%s\n", p.Print())
		})
	}
}

func BenchmarkAllValid(b *testing.B) {
	for _, test := range testcases {
		test := test
		if !test.readable || !test.solvable {
			continue
		}
		b.Run(test.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				var p sudoku.Puzzle
				p.ReadInput([]byte(test.in))
				p.Solve()
			}
		})
	}
}
