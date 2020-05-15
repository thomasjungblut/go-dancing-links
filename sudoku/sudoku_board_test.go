package sudoku

import (
	"github.com/stretchr/testify/assert"
	"os"
	"strings"
	"testing"
)

func TestReadingHappyPath(t *testing.T) {
	board := firstEulerGrid(t)
	builder := &strings.Builder{}
	assert.Nil(t, board.Print(builder))

	assert.Equal(t, `[0 0 3 0 2 0 6 0 0]
[9 0 0 3 0 5 0 0 1]
[0 0 1 8 0 6 4 0 0]
[0 0 8 1 0 2 9 0 0]
[7 0 0 0 0 0 0 0 8]
[0 0 6 7 0 8 2 0 0]
[0 0 2 6 0 9 5 0 0]
[8 0 0 2 0 3 0 0 9]
[0 0 5 0 1 0 3 0 0]`, strings.TrimSuffix(builder.String(), "\n"))

}

func TestSolvingHappyPath(t *testing.T) {
	board := firstEulerGrid(t)
	board, err := board.FindSingleSolution()
	assert.Nil(t, err)
	assert.Nil(t, board.VerifyCorrectness())
	assert.Nil(t, board.Print(os.Stdout))
}

func TestSudokuCorrectnessFailsRowConstraint(t *testing.T) {
	board := NewSudokuBoard(9)
	assert.Nil(t, board.ReadEulerTextFormat(`Grid00
483921657
967345821
251876493
548122976
729564138
136798245
372689514
814253769
695417382`))
	assert.EqualError(t, board.VerifyCorrectness(), "error in row 3: unique constraint violated: [5 4 8 1 2 2 9 7 6]")
}

func TestCheckDuplicates(t *testing.T) {
	assert.Nil(t, checkForDuplicates([]int{1, 2, 3}, 3))
	assert.EqualError(t, checkForDuplicates([]int{1, 2, 2}, 3), "unique constraint violated: [1 2 2]")
}

func firstEulerGrid(t *testing.T) SudokuBoardI {
	board := NewSudokuBoard(9)
	assert.Nil(t, board.ReadEulerTextFormat(`Grid 01
003020600
900305001
001806400
008102900
700000008
006708200
002609500
800203009
005010300`))
	return board
}
