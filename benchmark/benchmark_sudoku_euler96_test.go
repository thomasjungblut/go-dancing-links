package benchmark

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/thomasjungblut/go-dancing-links/sudoku"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"testing"
	"time"
)

func BenchmarkSolvingLongestTakingEulerPuzzle(t *testing.B) {
	eulerBoards, err := readAllEulerBoards()
	assert.Nil(t, err)

	// only run the one board that takes so long
	theBoard := eulerBoards[46]
	for i := 0; i < t.N; i++ {
		start := time.Now()
		board, err := theBoard.FindSingleSolution()
		assert.Nil(t, err)
		elapsed := time.Since(start)
		fmt.Println(fmt.Sprintf("solving the board took %s", elapsed))
		assert.Nil(t, board.VerifyCorrectness())
		_ = board.Print(os.Stdout)
	}
}

func BenchmarkSolvingAllEulerPuzzlesHappyPath(t *testing.B) {
	eulerBoards, err := readAllEulerBoards()
	assert.Nil(t, err)

	wg := sync.WaitGroup{}
	wg.Add(len(eulerBoards))

	for i := range eulerBoards {
		go func(i int, boardI sudoku.SudokuBoardI) {
			defer wg.Done()
			start := time.Now()
			boardResult, err := boardI.FindSingleSolution()
			elapsed := time.Since(start)
			assert.Nil(t, err)
			fmt.Println(fmt.Sprintf("solving board %d took %s", i, elapsed))
			assert.Nil(t, boardResult.VerifyCorrectness())
		}(i, eulerBoards[i])
	}

	wg.Wait()
}

func readAllEulerBoards() ([]sudoku.SudokuBoardI, error) {
	txt, err := ioutil.ReadFile("p096_sudoku.txt")
	if err != nil {
		return nil, err
	}

	var boards []sudoku.SudokuBoardI
	// bit of a hacky parser, but does the job pretty well
	for _, grid := range strings.Split(string(txt), "Grid") {
		if len(grid) == 0 {
			continue
		}
		board := sudoku.NewSudokuBoard(9)
		err = board.ReadEulerTextFormat(grid)
		if err != nil {
			return nil, err
		}
		boards = append(boards, board)
	}

	return boards, nil
}
