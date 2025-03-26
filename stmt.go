package mysql

import (
	"database/sql"
)

type implStmt struct {
	Stmt
	common
	sqlStmt *sql.Stmt
}

// Stmt is a prepared statement.
// A Stmt is safe for concurrent use by multiple goroutines.
type Stmt interface {
	// Close closes the statement.
	Close() error

	// Exec executes a prepared statement with the given arguments and
	// returns a Result summarizing the effect of the statement.
	Exec(args ...interface{}) (Result, error)

	// Query executes a prepared query statement with the given arguments
	// and returns the query results as a *Rows.
	Query(args ...interface{}) (Rows, error)

	// QueryRow executes a prepared query statement with the given arguments.
	// If an error occurs during the execution of the statement, that error will
	// be returned by a call to Scan on the returned *Row, which is always non-nil.
	// If the query selects no rows, the *Row's Scan will return ErrNoRows.
	// Otherwise, the *Row's Scan scans the first selected row and discards
	// the rest.
	//
	// Example usage:
	//
	//  var name string
	//  err := nameByUseridStmt.QueryRow(id).Scan(&name)
	QueryRow(args ...interface{}) (Row, error)

	// get sql stmt
	SQLStmt() *sql.Stmt
}

func (stmt *implStmt) SQLStmt() *sql.Stmt {
	return stmt.sqlStmt
}

func (stmt *implStmt) Close() error {
	return stmt.do(func() error {
		return stmt.sqlStmt.Close()
	})
}

func (stmt *implStmt) QueryRow(args ...interface{}) (row Row, err error) {
	err = stmt.do(func() error {
		var sqlRow *sql.Row
		sqlRow = stmt.sqlStmt.QueryRowContext(stmt.ctx, args...)
		row = &implRow{
			sqlRow: sqlRow,
			common: common{
				ctx:   stmt.ctx,
				retry: stmt.retry,
			},
		}
		return nil
	})
	return
}

func (stmt *implStmt) Query(args ...interface{}) (rows Rows, err error) {
	err = stmt.do(func() error {
		var sqlRows *sql.Rows
		sqlRows, err = stmt.sqlStmt.QueryContext(stmt.ctx, args...)
		rows = &implRows{
			sqlRows: sqlRows,
			common: common{
				ctx:   stmt.ctx,
				retry: stmt.retry,
			},
		}
		return err
	})
	return
}

func (stmt *implStmt) Exec(args ...interface{}) (result Result, err error) {
	err = stmt.do(func() error {
		var sqlResult sql.Result
		sqlResult, err = stmt.sqlStmt.ExecContext(stmt.ctx, args...)

		result = &implResult{
			sqlResult: sqlResult,
			common: common{
				ctx:   stmt.ctx,
				retry: stmt.retry,
			},
		}
		return err
	})
	return
}
