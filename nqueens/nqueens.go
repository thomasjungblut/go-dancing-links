package nqueens

import (
	"errors"
	"fmt"
	"github.com/thomasjungblut/go-dancing-links/dlx"
	"io"
	"regexp"
	"strconv"
)

type NQueensBoardI interface {
	// prints the given board row-wise as 'o' for empty space and 'x' for a placed queen.
	Print(writer io.StringWriter) error
	// verifies that the n-queens problem in this board is correctly solved, will return an error otherwise
	VerifyCorrectness() error
	// solves the n-queens problem with DLX and returns all of the solutions or an error.
	FindAllSolutions() ([]NQueensBoardI, error)
	// solves the n-queens problem with DLX and returns the count of the solutions
	CountAllSolutions() (int, error)
	// returns the given board as a dense two dimensional array, where true denotes a placed queen
	AsTwoDimArray() [][]bool
}

type placementCoordinate struct {
	row, col int
}

type NQueensBoard struct {
	placements map[placementCoordinate]bool
	n          int
}

func (b *NQueensBoard) AsTwoDimArray() [][]bool {
	board := make([][]bool, b.n, b.n)
	for i := 0; i < b.n; i++ {
		board[i] = make([]bool, b.n)
	}
	for p := range b.placements {
		board[p.row][p.col] = true
	}

	return board
}

func (b *NQueensBoard) Print(writer io.StringWriter) error {
	for row := 0; row < b.n; row++ {
		for col := 0; col < b.n; col++ {
			c := "o "
			if b.placements[placementCoordinate{row: row, col: col}] {
				c = "x "
			}
			_, err := writer.WriteString(c)
			if err != nil {
				return err
			}
		}
		_, err := writer.WriteString("\n")
		if err != nil {
			return err
		}
	}

	return nil
}

func (b *NQueensBoard) VerifyCorrectness() error {
	// this is quite inefficient for larger board sizes, but is not on the critical path anywhere I hope
	// check board sizes first
	expectedLen := len(b.placements)
	if expectedLen != b.n {
		return errors.New(fmt.Sprintf("unexpected number of queens: size is %d, but defined length is %d", expectedLen, b.n))
	}

	board := b.AsTwoDimArray()

	for i, row := range board {
		if expectedLen != len(row) {
			return errors.New(fmt.Sprintf("unexpected row length: in %dth array is %d, but defined length is %d", i, len(row), b.n))
		}
	}

	// iterate all queens in the array and search the row/col and diagonals from there
	queenCount := 0
	for r := 0; r < b.n; r++ {
		for c := 0; c < b.n; c++ {
			if board[r][c] {
				queenCount++
				// scan the col
				for rx := 0; rx < b.n; rx++ {
					if rx != r && board[rx][c] {
						return errors.New(fmt.Sprintf("found col conflict at %d/%d with %d/%d", r, c, rx, c))
					}
				}

				// scan the row
				for cx := 0; cx < b.n; cx++ {
					if cx != c && board[r][cx] {
						return errors.New(fmt.Sprintf("found row conflict at %d/%d with %d/%d", r, c, r, cx))
					}
				}

				// keep in mind: we omit both of the upwards (left/right) diagonal checks
				// since we implicitly gain them from iterating from the top left to bottom right

				// scan the diagonals left downwards
				rx := r
				cx := c
				for rx < b.n && cx >= 0 {
					if cx != c && rx != r && board[rx][cx] {
						return errors.New(fmt.Sprintf("found diagonal conflict at %d/%d with %d/%d", r, c, rx, cx))
					}
					rx++
					cx--
				}

				// scan the diagonal right downwards
				rx = r
				cx = c
				for rx < b.n && cx < b.n {
					if cx != c && rx != r && board[rx][cx] {
						return errors.New(fmt.Sprintf("found diagonal conflict at %d/%d with %d/%d", r, c, rx, cx))
					}
					rx++
					cx++
				}

			}
		}
	}

	if queenCount != b.n {
		return errors.New(fmt.Sprintf("unexpected number of queens: found %d, but defined length is %d", queenCount, b.n))
	}

	return nil
}

func (b *NQueensBoard) solve() [][]string {
	mat := dlx.NewDancingLinkMatrix()

	// add the row and col constraints
	for i := 0; i < b.n; i++ {
		mat.AppendColumn(fmt.Sprintf("r_%d", i))
	}

	for i := 0; i < b.n; i++ {
		mat.AppendColumn(fmt.Sprintf("c_%d", i))
	}

	// left bottom to top right diag
	for i := 0; i < 2*b.n-1; i++ {
		mat.AppendSecondaryColumn(fmt.Sprintf("d_%d", i))
	}
	// reversed
	for i := 0; i < 2*b.n-1; i++ {
		mat.AppendSecondaryColumn(fmt.Sprintf("rd_%d", i))
	}

	numConstraints := len(mat.Columns())
	// to fill the rows with the respective queen positions, we can devise some simple coordinate math:
	// the row constraint equals x
	// the column constraint equals N + y
	// for the diagonal constraint equals 2*N + (x + y)
	// for the reverse diagonal constraint equals (4*N-1) + (N – x + y - 1)
	for r := 0; r < b.n; r++ {
		for c := 0; c < b.n; c++ {
			constraint := make([]bool, numConstraints)
			constraint[r] = true
			constraint[b.n+c] = true

			constraint[2*b.n+r+c] = true
			constraint[(4*b.n-1)+(b.n-r+c-1)] = true

			_ = mat.AppendRow(fmt.Sprintf("queen_%d_%d", r, c), constraint)
		}
	}

	return mat.Solve()
}

func (b *NQueensBoard) CountAllSolutions() (int, error) {
	return len(b.solve()), nil
}

func (b *NQueensBoard) FindAllSolutions() ([]NQueensBoardI, error) {
	solutions := b.solve()
	var resultBoards []NQueensBoardI
	regex := regexp.MustCompile(`queen_(\d+)_(\d+)`)
	for _, solution := range solutions {
		if len(solution) != b.n {
			return nil, errors.New(fmt.Sprintf("didn't expect %d queens on an %d board", len(solution), b.n))
		}
		resultBoard := &NQueensBoard{n: b.n, placements: map[placementCoordinate]bool{}}
		for _, s := range solution {
			subMatch := regex.FindStringSubmatch(s)
			if subMatch != nil {
				row, err := strconv.Atoi(subMatch[1])
				if err != nil {
					return nil, err
				}
				col, err := strconv.Atoi(subMatch[2])
				if err != nil {
					return nil, err
				}
				resultBoard.placements[placementCoordinate{row: row, col: col}] = true
			}
		}
		resultBoards = append(resultBoards, resultBoard)
	}

	return resultBoards, nil
}

func NewNQueensBoard(n int) NQueensBoardI {
	return &NQueensBoard{n: n, placements: map[placementCoordinate]bool{}}
}

func newTestingNQueensBoard(a [][]bool) NQueensBoardI {
	placements := map[placementCoordinate]bool{}
	for r := 0; r < len(a); r++ {
		for c := 0; c < len(a[r]); c++ {
			if a[r][c] {
				placements[placementCoordinate{row: r, col: c}] = true
			}
		}
	}

	return &NQueensBoard{n: len(a), placements: placements}
}
