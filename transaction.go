package mysql

import (
	"database/sql"
)

type implTx struct {
	Tx
	common
	sqlTx *sql.Tx
}

// Tx is an in-progress database transaction.
//
// A transaction must end with a call to Commit or Rollback.
// After a call to Commit or Rollback, all operations on the transaction fail with ErrTxDone.
//
// The statements prepared for a transaction by calling the transaction's Prepare or Stmt methods are closed by the call to Commit or Rollback.
type Tx interface {

	// Commit commits the transaction.
	Commit() error

	// Rollback aborts the transaction.
	Rollback() error

	// Prepare creates a prepared statement for use within a transaction.
	//
	// The returned statement operates within the transaction and can no longer
	// be used once the transaction has been committed or rolled back.
	//
	// To use an existing prepared statement on this transaction, see Tx.Stmt.
	Prepare(query string) (Stmt, error)

	// Exec executes a query that doesn't return rows.
	// For example: an INSERT and UPDATE.
	Exec(query string, args ...interface{}) (Result, error)

	// Query executes a query that returns rows, typically a SELECT.
	Query(query string, args ...interface{}) (Rows, error)

	// QueryRow executes a query that is expected to return at most one row.
	// QueryRow always returns a non-nil value. Errors are deferred until
	// Row's Scan method is called.
	QueryRow(query string, args ...interface{}) (Row, error)

	// Stmt returns a transaction-specific prepared statement from
	// an existing statement.
	//
	// Example:
	//  updateMoney, err := db.Prepare("UPDATE balance SET money=money+? WHERE id=?")
	//  ...
	//  tx, err := db.Begin()
	//  ...
	//  res, err := tx.Stmt(updateMoney).Exec(123.45, 98293203)
	//
	// The returned statement operates within the transaction and will be closed
	// when the transaction has been committed or rolled back.
	Stmt(stmt Stmt) (Stmt, error)
}

func (tx *implTx) Query(query string, args ...interface{}) (rows Rows, err error) {
	err = tx.do(func() error {
		var sqlRows *sql.Rows
		if len(args) > 0 {
			sqlRows, err = tx.sqlTx.QueryContext(tx.ctx, query, args...)
		} else {
			sqlRows, err = tx.sqlTx.QueryContext(tx.ctx, query)
		}

		rows = &implRows{
			sqlRows: sqlRows,
			common: common{
				ctx:   tx.ctx,
				retry: tx.retry,
			},
		}

		return err
	})
	return
}

func (tx *implTx) QueryRow(query string, args ...interface{}) (row Row, err error) {
	err = tx.do(func() error {
		var sqlRow *sql.Row
		if len(args) > 0 {
			sqlRow = tx.sqlTx.QueryRowContext(tx.ctx, query, args...)
		} else {
			sqlRow = tx.sqlTx.QueryRowContext(tx.ctx, query)
		}

		row = &implRow{
			sqlRow: sqlRow,
			common: common{
				ctx:   tx.ctx,
				retry: tx.retry,
			},
		}

		return nil
	})
	return
}

func (tx *implTx) Prepare(query string) (stmt Stmt, err error) {
	err = tx.do(func() error {
		sqlStmt, err := tx.sqlTx.PrepareContext(tx.ctx, query)
		stmt = &implStmt{
			sqlStmt: sqlStmt,
			common: common{
				ctx:   tx.ctx,
				retry: tx.retry,
			},
		}
		return err
	})
	return
}

func (tx *implTx) Rollback() (err error) {
	err = tx.do(func() error {
		return tx.sqlTx.Rollback()
	})
	return err
}

func (tx *implTx) Commit() (err error) {
	err = tx.do(func() error {
		return tx.sqlTx.Commit()
	})
	return err
}

func (tx *implTx) Exec(query string, args ...interface{}) (result Result, err error) {
	err = tx.do(func() error {
		var sqlResult sql.Result
		if len(args) > 0 {
			sqlResult, err = tx.sqlTx.ExecContext(tx.ctx, query, args...)
		} else {
			sqlResult, err = tx.sqlTx.ExecContext(tx.ctx, query)
		}

		result = &implResult{
			sqlResult: sqlResult,
			common: common{
				ctx:   tx.ctx,
				retry: tx.retry,
			},
		}

		return err
	})
	return
}

func (tx *implTx) Stmt(stmtparam Stmt) (resStmt Stmt, err error) {
	err = tx.do(func() error {
		stmt := tx.sqlTx.StmtContext(tx.ctx, stmtparam.SQLStmt())
		resStmt = &implStmt{
			sqlStmt: stmt,
			common: common{
				ctx:   tx.ctx,
				retry: tx.retry,
			},
		}
		return nil
	})

	return
}
