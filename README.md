# mysql的驱动和sql接口的封装
为了方便使用，本模块对go的sql包中的接口进行了封装，
主要封装了context的支持、重试、日志和统计等功能

---------------------------------------

## 使用方法示例

```go
package main

import (
	"context"
	"fmt"
	"time"

	"gitlab.nolibox.com/skyteam/go-toml"
	"gitlab.nolibox.com/skyteam/go-log"
	"gitlab.nolibox.com/skyteam/go-mysql"
)

type Config struct {
	Logconf   log.Config   `toml:"log"`
	Mysqlconf mysql.Config `toml:"mysql"`
}

var config Config
var factory *mysql.Factory

// the following context of configuration is only example
/*
[log]
file_path = "./log/all.log"
level = "DEBUG"
# unit: Mb
max_size_mb = 2048
max_backups = 5
# "text" or "json"
formatter = "text"
show_file_line = true

[mysql]
dsn = "test:test@(127.0.0.1:3306)/t_db"
retry = 2
/*
[log]
file_path = "./log/all.log"
level = "DEBUG"
# unit: Mb
max_size_mb = 2048
max_backups = 5
# "text" or "json"
formatter = "text"
show_file_line = true

[mysql]
dsn = "test:test@(127.0.0.1:3306)/t_db"
#重试测试，下面的配置表示一个数据库操作最大访问3次(原始访问一次，重试2次)
retry = 2
#数据库连接池的最大空闲连接数
db_conn_pool_max_idle = 100
#数据库连接池的最大可用连接数量
db_conn_pool_max_open = 1000
#一个数据库连接的最大生存周期(单位纳秒)
db_conn_pool_max_lifetime = 10000000
*/
const ConfFile = "config.ini"

func init() {
	// init the config, only need init once in any project
	if _, err := toml.DecodeFile(ConfFile, &config); err != nil {
		fmt.Printf("fail to read config.||err=%v||config=%v", err, ConfFile)
		return
	}
	log.Init(&config.Logconf)
	
    // init the sql Factory, only need init once in any project
	factory = mysql.NewFactory()
	if err := factory.Open(&config.Mysqlconf); err != nil {
		fmt.Printf("fail to get sql Factory instantce || err=%v", err)
		factory = nil
	}
}

// 当你确定当前的查询语句最多返回一条语句的情况下，使用本函数
// QueryRow会自动释放资源,无需手动处理，也没有提供相应的方法
func testQueryRow(factory *mysql.Factory, ctx context.Context) {
	row, err := factory.New(ctx).QueryRow("select * from t_table where id=1")
	if err != nil {
		println("testQueryRow:" + err.Error())
		return
	}
	var setName, localName string
	err = row.Scan(&setName, &localName)
	println("the result are :", setName, localName)
	if err != nil {
		println("testQueryRow:" + err.Error())
		return
	}
}

func main() {

	defer log.Close()
	defer factory.ClosePool()
	d := time.Now().Add(1 * time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), d)

	// Even though ctx will be expired, it is good practice to call its
	// cancelation function in any case. Failure to do so may keep the
	// context and its parent alive longer than necessary.
	defer cancel()

	if factory != nil {
		testQueryRow(factory, ctx)
	} else {
		println("Fail to get factory, I will die~")
	}
}

```
