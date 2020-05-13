package dlx

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOneByOneMatrixCreation(t *testing.T) {
	mat := NewDancingLinkMatrix()
	mat.AppendColumn("1")
	err := mat.AppendRow("A", []bool{true})
	assert.Nil(t, err)

	assert.Equal(t, []string{"1"}, mat.Columns())
	assert.Equal(t, [][]bool{{true}}, mat.AsDenseMatrix())
}

func TestSparsenessMultiColumn(t *testing.T) {
	mat := NewDancingLinkMatrix()
	mat.AppendColumn("1")
	mat.AppendColumn("2")
	err := mat.AppendRow("A", []bool{true, false})
	assert.Nil(t, err)
	err = mat.AppendRow("A", []bool{false, true})
	assert.Nil(t, err)

	assert.Equal(t, []string{"1", "2"}, mat.Columns())
	assert.Equal(t, [][]bool{{true, false}, {false, true}}, mat.AsDenseMatrix())
}

func TestWikipediaExampleDataCorrectnessAsDenseMatrix(t *testing.T) {
	mat := NewWikipediaExampleMatrix(t)
	assert.Equal(t, []string{"1", "2", "3", "4", "5", "6", "7"}, mat.Columns())
	assert.Equal(t, []string{"A", "B", "C", "D", "E", "F"}, mat.Rows())
	expected := [][]bool{
		{true, false, false, true, false, false, true},
		{true, false, false, true, false, false, false},
		{false, false, false, true, true, false, true},
		{false, false, true, false, true, true, false},
		{false, true, true, false, false, true, true},
		{false, true, false, false, false, false, true},
	}
	assert.Equal(t, expected, mat.AsDenseMatrix())
}

func TestRowThatMismatchesColumns(t *testing.T) {
	mat := NewDancingLinkMatrix()
	mat.AppendColumn("a")
	err := mat.AppendRow("A", []bool{true, true, true})
	assert.EqualError(t, err, "column mismatch: have only 1 columns registered, but got 3")
}

func TestWikipediaExampleCoverColumn(t *testing.T) {
	mat := NewWikipediaExampleMatrix(t)
	assert.Nil(t, mat.CoverColumn(3))
	assert.Equal(t, [][]bool{
		{false, false, false, true, false, false, false},
		{false, false, false, true, false, false, false},
		{false, false, false, true, false, false, false},
		{false, false, true, false, true, true, false},
		{false, true, true, false, false, true, true},
		{false, true, false, false, false, false, true},
	}, mat.AsDenseMatrix())
}

func TestWikipediaExampleCoverColumnAndUncover(t *testing.T) {
	mat := NewWikipediaExampleMatrix(t)
	assert.Equal(t, 7, mat.NumUncoveredColumns())
	assert.Nil(t, mat.CoverColumn(3))
	assert.Equal(t, [][]bool{
		{false, false, false, true, false, false, false},
		{false, false, false, true, false, false, false},
		{false, false, false, true, false, false, false},
		{false, false, true, false, true, true, false},
		{false, true, true, false, false, true, true},
		{false, true, false, false, false, false, true},
	}, mat.AsDenseMatrix())

	assert.Equal(t, 6, mat.NumUncoveredColumns())
	assert.Nil(t, mat.UncoverColumn(3))
	assert.Equal(t, [][]bool{
		{true, false, false, true, false, false, true},
		{true, false, false, true, false, false, false},
		{false, false, false, true, true, false, true},
		{false, false, true, false, true, true, false},
		{false, true, true, false, false, true, true},
		{false, true, false, false, false, false, true},
	}, mat.AsDenseMatrix())
	assert.Equal(t, 7, mat.NumUncoveredColumns())
}

func TestUncoveringNotCoveredFails(t *testing.T) {
	mat := NewWikipediaExampleMatrix(t)
	err := mat.UncoverColumn(1)
	assert.EqualError(t, err, "column at 1 has not been covered yet")
}

func TestCoveringWithOutOfBoundsIndexFails(t *testing.T) {
	mat := NewWikipediaExampleMatrix(t)
	err := mat.UncoverColumn(15)
	assert.EqualError(t, err, "column at index 15 does not exist")
	err = mat.UncoverColumn(-1)
	assert.EqualError(t, err, "column at index -1 does not exist")
	err = mat.UncoverColumn(8)
	assert.EqualError(t, err, "column at index 8 does not exist")
	err = mat.UncoverColumn(7)
	assert.EqualError(t, err, "column at index 7 does not exist")
	// this should work
	err = mat.CoverColumn(0)
	assert.Nil(t, err)
	err = mat.CoverColumn(6)
	assert.Nil(t, err)
}

func TestUncoveringTwiceFails(t *testing.T) {
	mat := NewWikipediaExampleMatrix(t)
	err := mat.CoverColumn(1)
	assert.Nil(t, err)
	err = mat.UncoverColumn(1)
	assert.Nil(t, err)
	err = mat.UncoverColumn(1)
	assert.EqualError(t, err, "column at 1 has not been covered yet")
}

func TestSolvingWikipediaExample(t *testing.T) {
	mat := NewWikipediaExampleMatrix(t)
	result := mat.Solve()
	assert.Equal(t, 1, len(result))
	// we don't really care about the ordering of the result, but it should contain B-D-F
	assert.ElementsMatch(t, []string{"B", "D", "F"}, result[0])
}

func TestSolvingKnuthPaperExample(t *testing.T) {
	mat := NewKnuthPaperExampleMatrix(t)
	result := mat.Solve()
	assert.Equal(t, 1, len(result))
	assert.ElementsMatch(t, []string{"1", "4", "5"}, result[0])
}

func TestSolvingMultiSolutionExample(t *testing.T) {
	mat := NewDancingLinkMatrix()
	for i := 0; i < 3; i++ {
		mat.AppendColumn(fmt.Sprintf("%d", i))
	}
	assert.Nil(t, mat.AppendRow("A", []bool{true, true, true}))
	assert.Nil(t, mat.AppendRow("B", []bool{true, false, true}))
	assert.Nil(t, mat.AppendRow("C", []bool{false, true, false}))
	assert.Nil(t, mat.AppendRow("D", []bool{true, true, false}))
	assert.Nil(t, mat.AppendRow("E", []bool{false, false, true}))

	result := mat.Solve()
	assert.Equal(t, 3, len(result))
	assert.ElementsMatch(t, []string{"A"}, result[0])
	assert.ElementsMatch(t, []string{"B", "C"}, result[1])
	assert.ElementsMatch(t, []string{"D", "E"}, result[2])
}

func TestReadMeExample(t *testing.T) {
	mat := NewDancingLinkMatrix()

	mat.AppendColumn("beer")
	mat.AppendColumn("nachos")
	mat.AppendColumn("sour cream")

	_ = mat.AppendRow("Jack", []bool{true, false, false})  // Jack can bring beer only
	_ = mat.AppendRow("Amanda", []bool{true, true, false}) // Amanda can bring beer and nachos
	_ = mat.AppendRow("Chris", []bool{false, false, true}) // Chris can only bring sour cream
	_ = mat.AppendRow("Jen", []bool{true, true, true})     // Jen can bring everything

	result := mat.Solve()
	assert.Equal(t, 2, len(result))
	assert.ElementsMatch(t, []string{"Amanda", "Chris"}, result[0])
	assert.ElementsMatch(t, []string{"Jen"}, result[1])
}

func NewWikipediaExampleMatrix(t *testing.T) DancingLinksMatrixI {
	mat := NewDancingLinkMatrix()
	for i := 1; i < 8; i++ {
		mat.AppendColumn(fmt.Sprintf("%d", i))
	}

	assert.Nil(t, mat.AppendRow("A", []bool{true, false, false, true, false, false, true}))
	assert.Nil(t, mat.AppendRow("B", []bool{true, false, false, true, false, false, false}))
	assert.Nil(t, mat.AppendRow("C", []bool{false, false, false, true, true, false, true}))
	assert.Nil(t, mat.AppendRow("D", []bool{false, false, true, false, true, true, false}, ))
	assert.Nil(t, mat.AppendRow("E", []bool{false, true, true, false, false, true, true}))
	assert.Nil(t, mat.AppendRow("F", []bool{false, true, false, false, false, false, true}))
	return mat
}

func NewKnuthPaperExampleMatrix(t *testing.T) DancingLinksMatrixI {
	mat := NewDancingLinkMatrix()
	for i := 1; i < 8; i++ {
		mat.AppendColumn(fmt.Sprintf("%d", i))
	}

	assert.Nil(t, mat.AppendRow("1", []bool{false, false, true, false, true, true, false}))
	assert.Nil(t, mat.AppendRow("2", []bool{true, false, false, true, false, false, true}))
	assert.Nil(t, mat.AppendRow("3", []bool{false, true, true, false, false, true, false}))
	assert.Nil(t, mat.AppendRow("4", []bool{true, false, false, true, false, false, false}, ))
	assert.Nil(t, mat.AppendRow("5", []bool{false, true, false, false, false, false, true}))
	assert.Nil(t, mat.AppendRow("6", []bool{false, false, false, true, true, false, true}))
	return mat
}
