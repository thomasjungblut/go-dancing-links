package sudoku

import (
	"fmt"
	"github.com/thomasjungblut/go-dancing-links/dlx"
	"io"
	"math"
	"regexp"
	"strconv"
	"strings"
)

var NoSolutionError = fmt.Errorf("board has no solution")

type SudokuBoardI interface {
	// dense format used in Euler #96, first line is the header, rest is a 9x9 matrix without spaces and newline row separation
	ReadEulerTextFormat(gridString string) error
	// prints the given board row-wise in Go's array format
	Print(writer io.StringWriter) error
	// solves the sudoku with DLX by filling all zeros, when there are multiple solutions it
	// will pick the first and return it as a new board. If there are no solution it will be nil and an NoSolutionError.
	FindSingleSolution() (SudokuBoardI, error)
	// solves the sudoku with DLX by filling all zeros, multiple solutions are returned as a slice of new boards
	// if there are no solution it will be nil and an NoSolutionError.
	FindAllSolutions() ([]SudokuBoardI, error)
	// verifies that the Sudoku in this board is correctly solved
	VerifyCorrectness() error
}

type SudokuBoard struct {
	board [][]int
	size  int
}

func (b *SudokuBoard) FindSingleSolution() (SudokuBoardI, error) {
	boards, err := b.FindAllSolutions()
	if err != nil {
		return nil, err
	}
	if len(boards) > 0 {
		// pick the first if any
		return boards[0], nil
	} else {
		return nil, NoSolutionError
	}
}

func (b *SudokuBoard) FindAllSolutions() ([]SudokuBoardI, error) {
	squareYSize := int(math.Sqrt(float64(b.size)))
	squareXSize := b.size / squareYSize

	mat := dlx.NewDancingLinkMatrix()
	// column constraints
	for col := 0; col < b.size; col++ {
		for num := 1; num <= b.size; num++ {
			mat.AppendColumn(fmt.Sprintf("col_%d_%d", num, col))
		}
	}

	// row constraints
	for row := 0; row < b.size; row++ {
		for num := 1; num <= b.size; num++ {
			mat.AppendColumn(fmt.Sprintf("row_%d_%d", num, row))
		}
	}

	// square constraints
	for row := 0; row < squareXSize; row++ {
		for col := 0; col < squareYSize; col++ {
			for num := 1; num <= b.size; num++ {
				mat.AppendColumn(fmt.Sprintf("sq_%d_%d_%d", num, row, col))
			}
		}
	}

	// all cell constraints
	for row := 0; row < b.size; row++ {
		for col := 0; col < b.size; col++ {
			mat.AppendColumn(fmt.Sprintf("cell_%d_%d", row, col))
		}
	}

	// add existing board information
	for row := 0; row < b.size; row++ {
		for col := 0; col < b.size; col++ {
			if b.board[row][col] == 0 {
				// unknown cell, we have to add all constraints into the mix
				for num := 1; num <= b.size; num++ {
					err := mat.AppendRow(fmt.Sprintf("row_%d_%d_%d", row, col, num),
						b.generateRow(len(mat.Columns()), b.size, squareXSize, squareYSize, row, col, num))
					if err != nil {
						return nil, err
					}
				}
			} else {
				err := mat.AppendRow(fmt.Sprintf("row_%d_%d_%d", row, col, b.board[row][col]),
					b.generateRow(len(mat.Columns()), b.size, squareXSize, squareYSize, row, col, b.board[row][col]))
				if err != nil {
					return nil, err
				}
			}
		}
	}

	solutions := mat.Solve()
	if len(solutions) == 0 {
		return nil, NoSolutionError
	}

	var resultBoards []SudokuBoardI
	regex := regexp.MustCompile(`row_(\d)_(\d)_(\d)`)
	for _, solution := range solutions {
		board := make([][]int, b.size, b.size)
		for i := 0; i < b.size; i++ {
			board[i] = make([]int, b.size)
			copy(board[i], b.board[i])
		}
		resultBoard := &SudokuBoard{size: b.size, board: board}

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
				val, err := strconv.Atoi(subMatch[3])
				if err != nil {
					return nil, err
				}
				resultBoard.board[row][col] = val
			}
		}
		resultBoards = append(resultBoards, resultBoard)
	}

	return resultBoards, nil
}

func (b *SudokuBoard) generateRow(numCols, size, squareXSize, squareYSize, x, y, num int) []bool {
	r := make([]bool, numCols)
	xBox := int(x / squareXSize)
	yBox := int(y / squareYSize)
	// column
	r[x*size+num-1] = true
	// row
	r[size*size+y*size+num-1] = true
	// square
	r[2*size*size+(xBox*squareXSize+yBox)*size+num-1] = true
	// cell value
	r[3*size*size+size*x+y] = true
	return r
}

func (b *SudokuBoard) VerifyCorrectness() error {
	// check the rows for uniqueness
	for i := 0; i < b.size; i++ {
		err := checkForDuplicates(b.board[i], b.size)
		if err != nil {
			return fmt.Errorf("error in row %d: %v", i, err)
		}
	}

	// check the columns
	for i := 0; i < b.size; i++ {
		col := make([]int, b.size)
		for r := 0; r < b.size; r++ {
			col[r] = b.board[i][r]
		}
		err := checkForDuplicates(col, b.size)
		if err != nil {
			return fmt.Errorf("error in col %d: %v", i, err)
		}
	}

	// check the squares
	squareYSize := int(math.Sqrt(float64(b.size)))
	squareXSize := b.size / squareYSize

	for startRow := 0; startRow < b.size; startRow += squareXSize {
		for startCol := 0; startCol < b.size; startCol += squareYSize {

			idx := 0
			square := make([]int, b.size)
			for row := 0; row < squareXSize; row++ {
				for col := 0; col < squareYSize; col++ {
					square[idx] = b.board[row+startRow][col+startCol]
					idx++
				}
			}
			err := checkForDuplicates(square, b.size)
			if err != nil {
				return fmt.Errorf("error in square starting at %d/%d: %v", startRow, startCol, err)
			}
		}
	}

	return nil
}

func checkForDuplicates(numbers []int, size int) error {
	uniqueConstraint := make([]bool, size+1) // +1 since we're zero-indexing
	for _, n := range numbers {
		if uniqueConstraint[n] {
			return fmt.Errorf("unique constraint violated: %v", numbers)
		} else {
			uniqueConstraint[n] = true
		}
	}
	return nil
}

func (b *SudokuBoard) ReadEulerTextFormat(gridString string) error {
	// TODO(thomas): missing some validations
	for i, line := range strings.Split(strings.Trim(gridString, "\n"), "\n") {
		// ignore the first line, since it contains only the grid header
		if i > 0 {
			for j, c := range line {
				n, err := strconv.Atoi(string(c))
				if err != nil {
					return err
				}
				b.board[i-1][j] = n
			}
		}
	}
	return nil
}

func (b *SudokuBoard) Print(writer io.StringWriter) error {
	for _, row := range b.board {
		_, err := writer.WriteString(fmt.Sprintf("%v\n", row))
		if err != nil {
			return err
		}
	}

	return nil
}

func NewSudokuBoard(size int) SudokuBoardI {
	board := make([][]int, size, size)
	for i := 0; i < size; i++ {
		board[i] = make([]int, size)
	}
	return &SudokuBoard{size: size, board: board}
}
