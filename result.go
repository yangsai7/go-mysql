package mysql

import (
	"database/sql"
)

type implResult struct {
	Result
	common
	sqlResult sql.Result
}

// A Result summarizes an executed SQL command.
type Result interface {
	// LastInsertId returns the integer generated by the database
	// in response to a command. Typically this will be from an
	// "auto increment" column when inserting a new row. Not all
	// databases support this feature, and the syntax of such
	// statements varies.
	LastInsertId() (int64, error)

	// RowsAffected returns the number of rows affected by an
	// update, insert, or delete. Not every database or database
	// driver may support this.
	RowsAffected() (int64, error)
}

func (result *implResult) LastInsertId() (id int64, err error) {
	err = result.do(func() error {
		id, err = result.sqlResult.LastInsertId()
		return err
	})
	return
}

func (result *implResult) RowsAffected() (num int64, err error) {
	err = result.do(func() error {
		num, err = result.sqlResult.RowsAffected()
		return err
	})
	return
}
