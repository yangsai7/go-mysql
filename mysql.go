package mysql

import (
	"context"
	"database/sql"

	"github.com/go-sql-driver/mysql"
)

type NullTime = mysql.NullTime

type common struct {
	ctx   context.Context
	retry int
}

// Mysql is proxy for mysql database instance;
type implMysql struct {
	Mysql
	common
	db *sql.DB
}

// as AOP function, do log 、retry 、check timeout ...and so on
func (cm *common) do(action func() error) (err error) {
	retry := cm.retry
	if retry < 0 {
		retry = 0
	}
	for i := 0; i <= retry; i++ {
		select {
		case <-cm.ctx.Done():
			return cm.ctx.Err()
		default:
		}
		err = action()
		if err == nil {
			return
		}
	}
	return
}

// Mysql is a database handle representing a pool of zero or more underlying connections.
// It's safe for concurrent use by multiple goroutines.
type Mysql interface {
	// Begin starts a transaction. The default isolation level is dependent on
	// the driver.
	// The provided TxOptions is optional and may be nil if defaults should be used.
	// If a non-default isolation level is used that the driver doesn't support,
	// an error will be returned.
	// the function call func (db *DB) BeginTx(ctx context.Context, opts *TxOptions)
	Begin(opts *sql.TxOptions) (Tx, error)

	// Close closes the database, releasing any open resources.
	//
	// It is rare to Close a DB, as the DB handle is meant to be
	// long-lived and shared between many goroutines.
	Close() error

	// Driver returns the database's underlying driver.
	//Driver() driver.Driver

	// Exec executes a query without returning any rows.
	// The args are for any placeholder parameters in the query.
	// the function call sql.ExecContext(ctx context.Context, query string, args ...interface{})
	// TODO
	Exec(query string, args ...interface{}) (Result, error)

	// Ping verifies a connection to the database is still alive,
	// establishing a connection if necessary.
	Ping() error

	// Prepare creates a prepared statement for later queries or executions.
	// Multiple queries or executions may be run concurrently from the
	// returned statement.
	// The caller must call the statement's Close method
	// when the statement is no longer needed.
	Prepare(query string) (Stmt, error)

	// Query executes a query that returns rows, typically a SELECT.
	// The args are for any placeholder parameters in the query.
	Query(query string, args ...interface{}) (Rows, error)

	// QueryRow executes a query that is expected to return at most one row.
	// QueryRow always returns a non-nil value. Errors are deferred until
	// Row's Scan method is called.
	QueryRow(query string, args ...interface{}) (Row, error)

	// Stats returns database statistics.
	Stats() (sql.DBStats, error)
}

func (mysql *implMysql) Stats() (stat sql.DBStats, err error) {
	err = mysql.do(func() error {
		stat = mysql.db.Stats()
		return nil
	})
	return
}

func (mysql *implMysql) QueryRow(query string, args ...interface{}) (row Row, err error) {
	err = mysql.do(func() error {
		var sqlRow *sql.Row
		if len(args) > 0 {
			sqlRow = mysql.db.QueryRowContext(mysql.ctx, query, args...)
		} else {
			sqlRow = mysql.db.QueryRowContext(mysql.ctx, query)
		}
		row = &implRow{
			sqlRow: sqlRow,
			common: common{
				ctx:   mysql.ctx,
				retry: mysql.retry,
			},
		}
		return nil
	})
	return
}

func (mysql *implMysql) Prepare(query string) (stmt Stmt, err error) {
	err = mysql.do(func() error {
		sqlStmt, err := mysql.db.PrepareContext(mysql.ctx, query)
		stmt = &implStmt{
			sqlStmt: sqlStmt,
			common: common{
				ctx:   mysql.ctx,
				retry: mysql.retry,
			},
		}
		return err
	})
	return
}

func (mysql *implMysql) Ping() (err error) {
	err = mysql.do(func() error {
		err = mysql.db.Ping()
		return err
	})
	return err
}

func (mysql *implMysql) Close() (err error) {
	err = mysql.do(func() error {
		err = mysql.db.Close()
		return err
	})
	return err
}

func (mysql *implMysql) Begin(opts *sql.TxOptions) (tx Tx, err error) {
	err = mysql.do(func() error {
		sqlTx, err := mysql.db.BeginTx(mysql.ctx, opts)
		tx = &implTx{
			sqlTx: sqlTx,
			common: common{
				ctx:   mysql.ctx,
				retry: mysql.retry,
			},
		}
		return err
	})
	return
}

func (mysql *implMysql) Exec(query string, args ...interface{}) (result Result, err error) {
	err = mysql.do(func() error {
		var sqles sql.Result
		if len(args) > 0 {
			sqles, err = mysql.db.ExecContext(mysql.ctx, query, args...)
		} else {
			sqles, err = mysql.db.ExecContext(mysql.ctx, query)
		}

		result = &implResult{
			sqlResult: sqles,
			common: common{
				ctx:   mysql.ctx,
				retry: mysql.retry,
			},
		}

		return err
	})
	return
}

func (mysql *implMysql) Query(query string, args ...interface{}) (rows Rows, err error) {
	err = mysql.do(func() error {
		var sqlRows *sql.Rows

		if len(args) > 0 {
			sqlRows, err = mysql.db.QueryContext(mysql.ctx, query, args...)
		} else {
			sqlRows, err = mysql.db.QueryContext(mysql.ctx, query)
		}

		rows = &implRows{
			sqlRows: sqlRows,
			common: common{
				ctx:   mysql.ctx,
				retry: mysql.retry,
			},
		}

		return err
	})
	return
}
