[![Build Status](https://travis-ci.org/thomasjungblut/go-dancing-links.svg?branch=master)](https://travis-ci.org/thomasjungblut/go-dancing-links)

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

In this simple example, we have to solutions: either Jen brings everything or Amanda and Chris can bring their stuff individually. Let's see what DLX thinks about it:

```go

result := mat.Solve()
fmt.Println(result)
// [[Amanda Chris] [Jen]]
```  

Awesome! The result is a two dimensional slice of row names, because there can be multiple solutions for any given matrix. 

## Sudoku Solver


## N-Queens Solver



