package sudoku

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"testing"
)

func BenchmarkSolvingAllEulerPuzzlesHappyPath(t *testing.B) {
	boards, err := readAllEulerBoards()
	assert.Nil(t, err)

	var wg sync.WaitGroup
	wg.Add(len(boards))

	for i, board := range boards {
		go func(i int, board SudokuBoardI) {
			defer wg.Done()
			assert.Nil(t, board.Solve())
			assert.Nil(t, board.VerifyCorrectness())
			fmt.Println(fmt.Sprintf("solution for board %d", i))
			assert.Nil(t, board.Print(os.Stdout))
		}(i, board)
	}

	wg.Wait()
}

func readAllEulerBoards() ([]SudokuBoardI, error) {
	txt, err := ioutil.ReadFile("p096_sudoku.txt")
	if err != nil {
		return nil, err
	}

	var boards []SudokuBoardI
	// bit of a hacky parser, but does the job pretty well
	for _, grid := range strings.Split(string(txt), "Grid") {
		if len(grid) == 0 {
			continue
		}
		board := NewSudokuBoard(9)
		err = board.ReadEulerTextFormat(grid)
		if err != nil {
			return nil, err
		}
		boards = append(boards, board)
	}

	return boards, nil
}
