package dlx

import (
	"fmt"
	"math"
)

type DancingLinksMatrix struct {
	columnCovered     []bool
	numNodesPerColumn []int
	columnIdentifiers []string
	rowIdentifiers    []string
	columnNodes       []*Node
	head              *Node // top-left corner "head" of the matrix
}

type Node struct {
	left   *Node
	right  *Node
	top    *Node
	bottom *Node

	// probably unnecessary overhead to store this on every node
	rowIndex int
	colIndex int
}

func (m *DancingLinksMatrix) AppendColumn(columnIdentifier string) {
	m.appendColumnInternally(columnIdentifier, true)
}

func (m *DancingLinksMatrix) AppendSecondaryColumn(columnIdentifier string) {
	m.appendColumnInternally(columnIdentifier, false)
}

func (m *DancingLinksMatrix) appendColumnInternally(columnIdentifier string, primary bool) {
	newCol := &Node{colIndex: len(m.columnIdentifiers)}
	newCol.top = newCol
	newCol.bottom = newCol

	if primary {
		// insert the new column node in between tail and head
		tail := m.head.left
		tail.right = newCol
		newCol.left = tail
		newCol.right = m.head
		m.head.left = newCol
	} else {
		newCol.left = newCol
		newCol.right = newCol
	}

	// make sure we track the column values properly
	m.columnIdentifiers = append(m.columnIdentifiers, columnIdentifier)
	m.columnNodes = append(m.columnNodes, newCol)
	m.columnCovered = append(m.columnCovered, false)
	m.numNodesPerColumn = append(m.numNodesPerColumn, 0)
}

func (m *DancingLinksMatrix) AppendRow(rowIdentifier string, rowValues []bool) error {
	if len(rowValues) != len(m.columnIdentifiers) {
		return fmt.Errorf("column mismatch: have only %d columns registered, but got %d",
			len(m.columnIdentifiers), len(rowValues))
	}

	numRows := len(m.rowIdentifiers)

	var last *Node
	for i := 0; i < len(rowValues); i++ {
		// since this models a sparse matrix, we're only interested in true values
		if rowValues[i] {
			colTop := m.columnNodes[i]
			m.numNodesPerColumn[i]++
			bottom := colTop.top
			node := &Node{top: bottom, bottom: colTop, colIndex: i, rowIndex: numRows}
			bottom.bottom = node
			colTop.top = node

			if last != nil {
				rowHead := last.right
				node.left = last
				node.right = rowHead
				last.right = node
				rowHead.left = node
			} else {
				node.left = node
				node.right = node
			}

			last = node
		}
	}

	m.rowIdentifiers = append(m.rowIdentifiers, rowIdentifier)
	return nil
}

func (m *DancingLinksMatrix) CoverColumn(columnIndex int) error {
	if columnIndex < 0 || columnIndex >= len(m.columnCovered) {
		return fmt.Errorf("column at index %d does not exist", columnIndex)
	}

	if m.columnCovered[columnIndex] {
		return fmt.Errorf("column at %d is already covered", columnIndex)
	}

	// cover the header
	header := m.columnNodes[columnIndex]
	header.left.right = header.right
	header.right.left = header.left

	// go down the columns and unlink the respective rows from their columns
	row := header.bottom
	for row != header {
		node := row.right
		for node != row {
			node.bottom.top = node.top
			node.top.bottom = node.bottom
			node = node.right
			m.numNodesPerColumn[node.colIndex]--
		}

		row = row.bottom
	}

	m.columnCovered[columnIndex] = true

	return nil
}

func (m *DancingLinksMatrix) UncoverColumn(columnIndex int) error {
	if columnIndex < 0 || columnIndex >= len(m.columnCovered) {
		return fmt.Errorf("column at index %d does not exist", columnIndex)
	}
	if !m.columnCovered[columnIndex] {
		return fmt.Errorf("column at %d has not been covered yet", columnIndex)
	}

	header := m.columnNodes[columnIndex]
	row := header.top
	for row != header {
		node := row.left
		for node != row {
			node.bottom.top = node
			node.top.bottom = node
			node = node.left
			m.numNodesPerColumn[node.colIndex]++
		}
		row = row.top
	}

	header.right.left = header
	header.left.right = header

	m.columnCovered[columnIndex] = false

	return nil
}

func (m *DancingLinksMatrix) Columns() []string {
	return m.columnIdentifiers
}

func (m *DancingLinksMatrix) Rows() []string {
	return m.rowIdentifiers
}

func (m *DancingLinksMatrix) NumUncoveredColumns() int {
	count := 0
	for _, covered := range m.columnCovered {
		if !covered {
			count++
		}
	}
	return count
}

func (m *DancingLinksMatrix) AsDenseMatrix() [][] bool {
	denseMatrix := make([][]bool, len(m.rowIdentifiers))
	for i := range denseMatrix {
		denseMatrix[i] = make([]bool, len(m.columnIdentifiers))
	}

	for _, n := range m.columnNodes {
		cur := n.bottom
		for cur != n {
			denseMatrix[cur.rowIndex][cur.colIndex] = true
			cur = cur.bottom
		}
	}

	return denseMatrix
}

func (m *DancingLinksMatrix) Solve() [][]string {
	// allocate an empty slice with 100 capacity to avoid enlarging it all the time
	searchResult := m.search(make([]int, 0, 100))
	// map the row indices back to their names
	c := make([][]string, len(searchResult))
	for i, row := range searchResult {
		c[i] = make([]string, len(row))
		for ji, j := range row {
			c[i][ji] = m.rowIdentifiers[j]
		}
	}

	return c
}

func (m *DancingLinksMatrix) search(partialSolution []int) [][]int {
	if m.head.right == m.head {
		// we have to copy here to not interfere with other recursion steps changing the partial solution slice
		c := make([]int, len(partialSolution), len(partialSolution))
		copy(c, partialSolution)
		return [][]int{c}
	} else {
		result := make([][]int, 0, 0)
		nextColumn := m.chooseNext(m.head.right)
		_ = m.CoverColumn(nextColumn.colIndex)
		row := nextColumn.bottom
		for row != nextColumn {
			// we're adding the next eligible column to the solution
			partialSolution = append(partialSolution, row.rowIndex)
			node := row.right
			// all other columns that are true in that row now need to be covered too
			for node != row {
				_ = m.CoverColumn(node.colIndex)
				node = node.right
			}

			// recurse and gather any sub-solutions
			for _, recResult := range m.search(partialSolution) {
				if len(recResult) > 0 {
					result = append(result, recResult)
				}
			}

			// revert the last covering for the next column iteration
			partialSolution = partialSolution[:len(partialSolution)-1]
			node = row.left
			// all other columns that are true in that row now need to be covered too
			for node != row {
				_ = m.UncoverColumn(node.colIndex)
				node = node.left
			}

			row = row.bottom
		}
		_ = m.UncoverColumn(nextColumn.colIndex)

		return result
	}
}

func (m *DancingLinksMatrix) chooseNext(node *Node) *Node {
	lowestCount := math.MaxInt32
	lowestNode := node
	for node != m.head {
		cnt := m.numNodesPerColumn[node.colIndex]
		if cnt < lowestCount {
			lowestNode = node
			lowestCount = cnt
		}
		node = node.right
	}
	return lowestNode
}

func NewDancingLinkMatrix() DancingLinksMatrixI {
	header := &Node{}
	header.left = header
	header.right = header
	header.top = header
	header.bottom = header

	return &DancingLinksMatrix{
		columnIdentifiers: []string{},
		rowIdentifiers:    []string{},
		columnNodes:       []*Node{},
		head:              header,
	}
}
