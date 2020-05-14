package benchmark

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/thomasjungblut/go-dancing-links/sudoku"
	"io/ioutil"
	"strings"
	"sync"
	"testing"
	"time"
)

func BenchmarkSolvingAllEulerPuzzlesHappyPath(t *testing.B) {
	eulerBoards, err := readAllEulerBoards()
	assert.Nil(t, err)

	wg := sync.WaitGroup{}
	wg.Add(len(eulerBoards))

	for i := range eulerBoards {
		go func(i int, boardI sudoku.SudokuBoardI) {
			defer wg.Done()
			start := time.Now()
			boardResults, err := boardI.FindAllSolutions()
			elapsed := time.Since(start)
			assert.Nil(t, err)
			fmt.Println(fmt.Sprintf("board %d has %d solution(s), took %s", i, len(boardResults), elapsed))
			for _, board := range boardResults {
				assert.Nil(t, board.VerifyCorrectness())
			}
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
