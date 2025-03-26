package mysql

import "time"

// Config 是 Mysql 配置信息。
type Config struct {
	Dsn         string        `toml:"dsn" yaml:"dsn"`                                // data source name
	Retry       int           `toml:"retry" yaml:"retry"`                            // retry time
	MaxIdle     int           `toml:"db_conn_pool_max_idle" yaml:"max_idle"`         // zero means defaultMaxIdleConns; negative means 0
	MaxOpen     int           `toml:"db_conn_pool_max_open" yaml:"max_open"`         // <= 0 means unlimited
	MaxLifetime time.Duration `toml:"db_conn_pool_max_lifetime" yaml:"max_lifetime"` // maximum amount of time a connection may be reused
}
