[![Go](https://github.com/thomasjungblut/go-dancing-links/actions/workflows/go.yml/badge.svg)](https://github.com/thomasjungblut/go-dancing-links/actions/workflows/go.yml)

## go-dancing-links

`go-dancing-links` contains a golang implementation of the dancing links algorithm (DLX) as Donald Knuth devised it in [his paper](https://arxiv.org/abs/cs/0011047).

The algorithm solves the [exact cover](https://en.wikipedia.org/wiki/Exact_cover) problem, which is a common constraint satisfaction algorithm that can be used to solve Sudokus or the n-queens problem.


## Installation

> go get -d github.com/thomasjungblut/go-dancing-links

## Using go-dancing-links

It's fairly easy to use, as long as you can translate your problem into a binary matrix.

Let's say you're at a (somewhat nerdy) party and you want to have exactly one of each item: beer, nachos and sour cream. Now your friends can bring some of them in varying configurations and you want to tell them whether they should bring their items or not. Fret not, DLX has you covered and in no time you'll have the solution(s) to your problem!

### Example

We need to tell the matrix what items you need (columns) and what friends you have and what they can bring to your party (rows).

```go

mat := NewDancingLinkMatrix()

mat.AppendColumn("beer")
mat.AppendColumn("nachos")
mat.AppendColumn("sour cream")

mat.AppendRow("Jack", []bool{true, false, false}) // Jack can bring beer only
mat.AppendRow("Amanda", []bool{true, true, false}) // Amanda can bring beer and nachos
mat.AppendRow("Chris", []bool{false, false, true}) // Chris can only bring sour cream
mat.AppendRow("Jen", []bool{true, true, true}) // Jen can bring everything 

``` 

In this simple example, we have two solutions: either Jen brings everything or Amanda and Chris can bring their stuff individually. Let's see what DLX thinks about it:

```go

result := mat.Solve()
fmt.Println(result)
// [[Amanda Chris] [Jen]]
```  

Awesome! The result is a two dimensional slice of row names, because there can be multiple solutions for any given matrix. 

## Sudoku Solver

Sudokus can also be solved pretty fast from a string by using the Euler96 format:

```go

board := NewSudokuBoard(9)
board.ReadEulerTextFormat(`Grid 01
    003020600
    900305001
    001806400
    008102900
    700000008
    006708200
    002609500
    800203009
    005010300`)
board, err := board.FindSingleSolution() // solve with DLX and return the filled board
err := board.VerifyCorrectness() // check if the solution is correct
err := board.Print(os.Stdout) // print it to stdout

// or if there are multiple solutions possible, you can find all boards
boards, err := board.FindAllSolutions() 
for _, board := range boards {
	err := board.VerifyCorrectness() // check if the solution is correct
	err := board.Print(os.Stdout) // print it to stdout
}
```

## N-Queens Solver

The generalized N-Queens problem can also be solved fairly easy with DLX:

```go

board := NewNQueensBoard(4)
result, err := board.FindAllSolutions()
for _, board := range result {
    assert.Nil(t, board.VerifyCorrectness())
    assert.Nil(t, board.Print(os.Stdout))
    println()
}

```

which outputs the results as a NxN block where 'x' denotes a queen:

```
o x o o 
o o o x 
x o o o 
o o x o 

o o x o 
x o o o 
o o o x 
o x o o 
```
