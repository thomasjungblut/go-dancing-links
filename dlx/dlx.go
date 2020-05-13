package dlx

type DancingLinksMatrixI interface {
	// Append a new column with the given name to the matrix
	AppendColumn(columnIdentifier string)
	// Append a given dense row to the matrix, error is returned when the number of columns mismatch the registered ones
	AppendRow(rowIdentifier string, rowValues []bool) error
	// Returns all column identifiers
	Columns() []string
	// Returns all row identifiers
	Rows() []string
	// Returns the internal doubly-linked-list structure as a dense matrix of booleans
	AsDenseMatrix() [][]bool

	// Covers the given column, meaning it will unlink the whole column and all the rows where the column is true.
	// error is returned when the column is already covered.
	CoverColumn(columnIndex int) error
	// Uncovers the column at the given index again, this undoes the CoverColumn operation.
	// error is returned when the given column has not been covered before.
	UncoverColumn(columnIndex int) error
	// Returns the number of not yet covered columns (uncovered columns)
	NumUncoveredColumns() int

	// Solves this matrix, returns the results as a list, of which each element is a set of rows that covers all the columns.
	// the first dimension would contain the number of solutions.
	// the second dimension contains the identifier of the rows that participate in this solution.
	Solve() [][]string
}
