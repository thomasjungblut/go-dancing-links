package benchmark

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/thomasjungblut/go-dancing-links/nqueens"
	"sync"
	"testing"
	"time"
)

func BenchmarkSolvingHighNumberNQueenProblems(t *testing.B) {
	expectedResultSizes := []int{
		1, 1, 0, 0, 2, 10, 4, 40, 92, 352, 724, 2680, 14200, 73712, 365596,
		2279184, 14772512, // 95815104, // 666090624, 4968057848, 39029188884,
	}

	wg := sync.WaitGroup{}
	wg.Add(len(expectedResultSizes))

	for i, expectedResultSize := range expectedResultSizes {
		go func(i, expectedResult int) {
			defer wg.Done()
			start := time.Now()
			c, err := nqueens.NewNQueensBoard(i).CountAllSolutions()
			elapsed := time.Since(start)
			assert.Nil(t, err)
			assert.Equal(t, expectedResult, c)
			fmt.Println(fmt.Sprintf("n = %d has %d solution(s), took %s", i, c, elapsed))
		}(i, expectedResultSize)
	}

	wg.Wait()
}
