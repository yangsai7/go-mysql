package mysql

import (
	"context"
	"database/sql"
	"log"

	"github.com/go-sql-driver/mysql" // 默认使用自带的 driver。
)

func init() {
	mysql.SetLogger(log.Default())
}

// Factory 存储了数据库的配置并可以创建 Mysql 客户端实例。
type Factory struct {
	config *Config
	db     *sql.DB
	mock   Mysql
}

// NewFactory 创建一个新的 Mysql 工厂，用New方法创建mysql实例。
// 通过被工厂方法产生的myslq实例
func NewFactory() (factory *Factory) {
	return &Factory{}
}

// Open 打开数据库连接和配置连接池
func (f *Factory) Open(config *Config) error {
	db, err := sql.Open("mysql", config.Dsn)
	if err != nil {
		log.Printf("Error while init mysql db with addr: %v, err: %v", config.Dsn, err)
		return err
	}
	db.SetMaxIdleConns(config.MaxIdle)
	db.SetMaxOpenConns(config.MaxOpen)
	if config.MaxLifetime != 0 {
		db.SetConnMaxLifetime(config.MaxLifetime)
	}

	if err = db.Ping(); err != nil {
		log.Printf("Error while ping mysql db with addr: %v, maxIdle: %v,"+
			"maxOpen: %v, idleMaxTime: %v, err: %v", config.Dsn, config.MaxIdle, config.MaxOpen, config.MaxLifetime, err)
		return err
	}

	f.config = config
	f.db = db
	return nil
}

// New 创建 Mysql 的客户端实例。
// 通过 Mysql 接口产生的 Rows 和 Stat 对象需要手动关闭释放连接资源。
// 通过 Tx 对象产生的 Rows 和 Stat 对象需要手动关闭释放连接资源。
// 其它的情况无需手动资源。
func (f *Factory) New(ctx context.Context) (mysql Mysql) {
	if f.mock != nil {
		return f.mock
	}

	return &implMysql{
		db: f.db,
		common: common{
			ctx:   ctx,
			retry: f.config.Retry,
		},
	}
}

// Close 关闭数据库连接池
// 除非你有确定的关闭理由，否则只需要在main方法中加上defer处理即可；
func (f *Factory) Close() error {
	return f.db.Close()
}

// SetMock 设置一个 mock 实例，所有的 New 都会返回这个 mock 实例。
func (f *Factory) SetMock(mock Mysql) {
	f.mock = mock
}

// ResetMock 清除 mock。
func (f *Factory) ResetMock() {
	f.SetMock(nil)
}
