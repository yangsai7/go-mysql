package mysql

import (
	"database/sql"
)

type implRow struct {
	Row
	common
	sqlRow *sql.Row
}

// Row is the result of calling QueryRow to select a single row.
type Row interface {
	// Scan copies the columns from the matched row into the values
	// pointed at by dest. See the documentation on Rows.Scan for details.
	// If more than one row matches the query,
	// Scan uses the first row and discards the rest. If no row matches
	// the query, Scan returns ErrNoRows.
	Scan(dest ...interface{}) error
}

func (row *implRow) Scan(dest ...interface{}) error {
	return row.do(func() error {
		return row.sqlRow.Scan(dest...)
	})
}
