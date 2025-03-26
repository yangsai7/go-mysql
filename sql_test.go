package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"testing"
	"time"
)

const fakeDBName = "foo"

var chrisBirthday = time.Unix(123456789, 0)

func mysqlExec(t testing.TB, mysql Mysql, query string, args ...interface{}) {
	result, err := mysql.Exec(query, args...)
	if err != nil {
		t.Fatalf("Exec of %q: %v", query, err)
	}
	affectRow, err := result.RowsAffected()
	t.Logf("query %v,RowsAffected %d", query, affectRow)

}

func getFactory(t *testing.T) *Factory {
	config := &Config{
		Retry: 2,
	}

	db, _ := sql.Open("test", fakeDBName)
	factory := &Factory{
		config: config,
		db:     db,
	}

	d := time.Now().Add(1 * time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), d)
	defer cancel()
	mysql := factory.New(ctx)
	if _, err := mysql.Exec("WIPE"); err != nil {
		t.Fatalf("exec wipe: %v", err)
	}
	mysqlExec(t, mysql, "CREATE|people|name=string,age=int32,photo=blob,dead=bool,bdate=datetime")
	mysqlExec(t, mysql, "INSERT|people|name=Alice,age=?,photo=APHOTO", 1)
	mysqlExec(t, mysql, "INSERT|people|name=Bob,age=?,photo=BPHOTO", 2)
	mysqlExec(t, mysql, "INSERT|people|name=Chris,age=?,photo=CPHOTO,bdate=?", 3, chrisBirthday)

	return factory
}

func TestQueryRow(t *testing.T) {
	d := time.Now().Add(1 * time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), d)
	defer cancel()
	mysql := getFactory(t).New(ctx)
	defer closeDB(t, mysql)
	var name string
	var age int
	var birthday time.Time

	row, err := mysql.QueryRow("SELECT|people|bdate|age=?|LIMIT", 3)

	if err != nil {
		t.Fatalf("fail to QueryRow. [err:%v]", err)
	}

	err = row.Scan(&birthday)
	if err != nil || !birthday.Equal(chrisBirthday) {
		t.Fatalf("chris birthday = %v, err = %v; want %v", birthday, err, chrisBirthday)
	}
	row, err = mysql.QueryRow("SELECT|people|age,name|age=?|LIMIT", 2)
	err = row.Scan(&age, &name)
	if err != nil {
		t.Fatalf("age QueryRow+Scan: %v", err)
	}

	if name != "Bob" {
		t.Fatalf("expected name Bob, got %q", name)
	}
	if age != 2 {
		t.Fatalf("expected age 2, got %d", age)
	}

	row, err = mysql.QueryRow("SELECT|people|age,name|name=?|limit", "Alice")
	err = row.Scan(&age, &name)
	if err != nil {
		t.Fatalf("name QueryRow+Scan: %v", err)
	}
	if name != "Alice" {
		t.Fatalf("expected name Alice, got %q", name)
	}
	if age != 1 {
		t.Fatalf("expected age 1, got %d", age)
	}

	var photo []byte
	row, err = mysql.QueryRow("SELECT|people|photo|name=?|limit", "Alice")
	err = row.Scan(&photo)
	if err != nil {
		t.Fatalf("photo QueryRow+Scan: %v", err)
	}
	want := []byte("APHOTO")
	if !reflect.DeepEqual(photo, want) {
		t.Fatalf("photo = %q; want %q", photo, want)
	}
}

func TestExec(t *testing.T) {
	d := time.Now().Add(1 * time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), d)
	defer cancel()
	mysql := getFactory(t).New(ctx)
	defer closeDB(t, mysql)
	mysqlExec(t, mysql, "CREATE|t1|name=string,age=int32,dead=bool")
	stmt, err := mysql.Prepare("INSERT|t1|name=?,age=?")
	if err != nil {
		t.Fatalf("Stmt, err = %v, %v", stmt, err)
	}
	defer stmt.Close()

	type execTest struct {
		args    []interface{}
		wantErr string
	}
	execTests := []execTest{
		// Okay:
		{[]interface{}{"Brad", 31}, ""},
		{[]interface{}{"Brad", int64(31)}, ""},
		{[]interface{}{"Bob", "32"}, ""},
		{[]interface{}{7, 9}, ""},

		// Invalid conversions:
		{[]interface{}{"Brad", int64(0xFFFFFFFF)}, "sql: converting argument $2 type: sql/driver: value 4294967295 overflows int32"},
		{[]interface{}{"Brad", "strconv fail"}, `sql: converting argument $2 type: sql/driver: value "strconv fail" can't be converted to int32`},

		// Wrong number of args:
		{[]interface{}{}, "sql: expected 2 arguments, got 0"},
		{[]interface{}{1, 2, 3}, "sql: expected 2 arguments, got 3"},
	}
	for n, et := range execTests {
		_, err := stmt.Exec(et.args...)
		errStr := ""
		if err != nil {
			errStr = err.Error()
		}
		if errStr != et.wantErr {
			t.Fatalf("stmt.Execute #%d: for %v, got error %q, want error %q",
				n, et.args, errStr, et.wantErr)
		}
	}
}

func TestTxPrepare(t *testing.T) {
	d := time.Now().Add(1 * time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), d)
	defer cancel()
	mysql := getFactory(t).New(ctx)
	defer closeDB(t, mysql)
	mysqlExec(t, mysql, "CREATE|t1|name=string,age=int32,dead=bool")
	tx, err := mysql.Begin(nil)
	if err != nil {
		t.Fatalf("Begin = %v", err)
	}
	stmt, err := tx.Prepare("INSERT|t1|name=?,age=?")
	if err != nil {
		t.Fatalf("Stmt, err = %v, %v", stmt, err)
	}
	defer stmt.Close()

	_, err = stmt.Exec("Bobby", 7)
	if err != nil {
		t.Fatalf("Exec = %v", err)
	}

	err = tx.Commit()
	if err != nil {
		t.Fatalf("Commit = %v", err)
	}
}

func TestTxStmt(t *testing.T) {
	d := time.Now().Add(1 * time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), d)
	defer cancel()
	mysql := getFactory(t).New(ctx)
	defer closeDB(t, mysql)
	mysqlExec(t, mysql, "CREATE|t1|name=string,age=int32,dead=bool")
	stmt, err := mysql.Prepare("INSERT|t1|name=?,age=?")
	if err != nil {
		t.Fatalf("Stmt, err = %v, %v", stmt, err)
	}
	defer stmt.Close()
	tx, err := mysql.Begin(nil)
	if err != nil {
		t.Fatalf("Begin = %v", err)
	}
	txs, _ := tx.Stmt(stmt)
	defer txs.Close()
	_, err = txs.Exec("Bobby", 7)
	if err != nil {
		t.Fatalf("Exec = %v", err)
	}
	err = tx.Commit()
	if err != nil {
		t.Fatalf("Commit = %v", err)
	}
}

func TestStmtQuery(t *testing.T) {
	d := time.Now().Add(1 * time.Second)
	ctx, _ := context.WithDeadline(context.Background(), d)
	mysql := getFactory(t).New(ctx)
	defer closeDB(t, mysql)
	mysqlExec(t, mysql, "CREATE|t1|name=string,age=int32,dead=bool")
	mysqlExec(t, mysql, "INSERT|t1|name=Alice")

	stmt, err := mysql.Prepare("SELECT|t1|name|name=?|limit")
	rows, err := stmt.Query("Alice")
	defer stmt.Close()
	defer rows.Close()
	if err != nil {
		t.Fatal(err)
	}
	has := rows.Next()
	for has {
		var name string
		rows.Scan(&name)
		t.Logf("query rows: name:%v", name)
		has = rows.Next()
	}

	stmt, _ = mysql.Prepare("SELECT|t1|name|name=?|limit")
	row, err := stmt.QueryRow("Alice")
	if err != nil {
		t.Fatal(err)
	}
	var name string
	row.Scan(&name)
	t.Logf("query row:name:%v", name)

}

func TestTxQuery(t *testing.T) {
	d := time.Now().Add(1 * time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), d)
	defer cancel()
	mysql := getFactory(t).New(ctx)
	defer closeDB(t, mysql)
	mysqlExec(t, mysql, "CREATE|t1|name=string,age=int32,dead=bool")
	mysqlExec(t, mysql, "INSERT|t1|name=Alice")

	tx, err := mysql.Begin(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()

	r, err := tx.Query("SELECT|t1|name||limit")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()

	if ok := r.Next(); !ok {
		if r.Err() != nil {
			t.Fatal(r.Err())
		}
		t.Fatal("expected one row")
	}

	var x string
	err = r.Scan(&x)
	if err != nil {
		t.Fatal(err)
	}
}

func TestTxExec(t *testing.T) {
	d := time.Now().Add(1 * time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), d)
	defer cancel()
	mysql := getFactory(t).New(ctx)
	defer closeDB(t, mysql)
	mysqlExec(t, mysql, "CREATE|t1|name=string,age=int32,dead=bool")

	tx, err := mysql.Begin(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()
	result, err := tx.Exec("INSERT|t1|name=Alice")
	if err != nil {
		t.Fatal(err)
	} else {
		affectRow, err := result.RowsAffected()
		if err == nil {
			t.Logf("TestTxExec %d", affectRow)
		} else {
			t.Fatal(err)
		}
	}
}

func closeDB(t testing.TB, mysql Mysql) {
	if e := recover(); e != nil {
		fmt.Printf("Panic: %v\n", e)
		panic(e)
	}
	defer setHookpostCloseConn(nil)
	setHookpostCloseConn(func(_ *fakeConn, err error) {
		if err != nil {
			t.Fatalf("Error closing fakeConn: %v", err)
		}
	})

	err := mysql.Close()
	if err != nil {
		t.Fatalf("error closing DB: %v", err)
	}

}
