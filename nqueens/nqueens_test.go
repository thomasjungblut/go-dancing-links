package nqueens

import (
	"github.com/stretchr/testify/assert"
	"os"
	"strings"
	"testing"
)

func TestFourQueensSolutionsReadMe(t *testing.T) {
	board := NewNQueensBoard(4)
	result, err := board.FindAllSolutions()
	assert.Equal(t, 2, len(result))
	assert.Nil(t, err)
	for _, board := range result {
		assert.Nil(t, board.VerifyCorrectness())
		assert.Nil(t, board.Print(os.Stdout))
		println()
	}
}

func TestNQueensResultSizesAndCheckCorrectness(t *testing.T) {
	// https://oeis.org/A000170
	expectedResultSizes := []int{
		1, 1, 0, 0, 2, 10, 4, 40, 92, 352, 724, 2680, 14200, 73712,
	}

	for i := 0; i < len(expectedResultSizes); i++ {
		board := NewNQueensBoard(i)
		result, err := board.FindAllSolutions()
		assert.Nil(t, err)
		assert.Equal(t, expectedResultSizes[i], len(result), "n = %d", i)
		for _, board := range result {
			assert.Nil(t, board.VerifyCorrectness(), "n = %d", i)
		}
	}
}

func TestPrintingHappyPath(t *testing.T) {
	board := newTestingNQueensBoard([][]bool{
		{true, false, false},
		{false, false, true},
		{false, false, false},
	})
	builder := &strings.Builder{}
	assert.Nil(t, board.Print(builder))

	assert.Equal(t, `x o o 
o o x 
o o o `, strings.TrimSuffix(builder.String(), "\n"))

}

func TestVerifyHappyPath(t *testing.T) {
	board := newTestingNQueensBoard([][]bool{
		{true},
	})
	assert.Nil(t, board.VerifyCorrectness())
}

func TestVerifyWrongSizesFails(t *testing.T) {
	board := newTestingNQueensBoard([][]bool{
		{true},
		{true, false},
	})
	assert.EqualError(t, board.VerifyCorrectness(), "unexpected row length: in 0th array is 1, but defined length is 2")
}

func TestVerifyWrongResultFailsDiagonal(t *testing.T) {
	board := newTestingNQueensBoard([][]bool{
		{true, false},
		{false, true},
	})
	assert.EqualError(t, board.VerifyCorrectness(), "found diagonal conflict at 0/0 with 1/1")
}

func TestVerifyWrongResultFailsDiagonalBackwards(t *testing.T) {
	board := newTestingNQueensBoard([][]bool{
		{false, true},
		{true, false},
	})
	assert.EqualError(t, board.VerifyCorrectness(), "found diagonal conflict at 0/1 with 1/0")
}

func TestVerifyWrongResultFailsRow(t *testing.T) {
	board := newTestingNQueensBoard([][]bool{
		{true, true},
		{false, true},
	})
	assert.EqualError(t, board.VerifyCorrectness(), "found row conflict at 0/0 with 0/1")
}

func TestVerifyWrongResultFailsColumn(t *testing.T) {
	board := newTestingNQueensBoard([][]bool{
		{true, false},
		{true, false},
	})
	assert.EqualError(t, board.VerifyCorrectness(), "found col conflict at 0/0 with 1/0")
}

func TestVerifyWrongResultFailsWrongQueenCount(t *testing.T) {
	board := newTestingNQueensBoard([][]bool{
		{true, false},
		{false, false},
	})
	assert.EqualError(t, board.VerifyCorrectness(), "unexpected number of queens: found 1, but defined length is 2")
}
