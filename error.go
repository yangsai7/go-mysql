package mysql

import (
	"fmt"

	"github.com/go-sql-driver/mysql" // 默认使用自带的 driver。
)

// Error 表示一个 Mysql 详细错误信息。
type Error struct {
	Number  int
	Message string
}

// Error 返回错误信息。
func (e *Error) Error() string {
	return fmt.Sprintf("error %v: %v", e.Number, e.Message)
}

// ParseError 尝试从 err 中解析 Mysql 错误信息，如果解析不出来则返回 nil。
func ParseError(err error) *Error {
	if err == nil {
		return nil
	}

	if e, ok := err.(*mysql.MySQLError); ok {
		return &Error{
			Number:  int(e.Number),
			Message: e.Message,
		}
	}

	return nil
}
