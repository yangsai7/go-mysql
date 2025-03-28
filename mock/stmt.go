// Autogenerated by the tool "go-mock-code-gen" built from
// gitlab.nolibox.com/skyteam/go-mock-code-gen
// DO NOT EDIT UNLESS YOU ARE SURE THAT YOU KNOW WHAT YOU ARE DOING

// mock 可以 mock 掉 Stmt interface 的所有方法，方便编写单元测试。
package mock

import (
	"database/sql"
	"fmt"
	"sync"

	mysql "github.com/yangsai7/go-mysql"
)

type mockDataStmt map[string][]interface{}
type mockDefaultDataStmt map[string]interface{}

// MockDataAdderStmt 用于向 MockStmt 添加指定接口的测试数据。
// 这里添加的测试数据是一次性的，只要设置了的 mock 数据的接口被调用，数据就会按照 FIFO 顺序消费掉。
type MockDataAdderStmt struct {
	m *MockStmt
}

// MockFuncAdderStmt 用于向 MockStmt 添加指定接口的测试方法，当接口调用时这里设置的方法就会调用。
// 这里添加的测试数据是一次性的，只要设置了的 mock 数据的接口被调用，数据就会按照 FIFO 顺序消费掉。
type MockFuncAdderStmt struct {
	m *MockStmt
}

// MockDefaultDataStmt 用于向 MockStmt 添加指定接口的默认测试数据。
// 如果该接口没有设置测试数据或者测试数据已经消费完，这个默认的测试数据就会被返回给调用者。
// MockDefaultDataStmt 和 MockDefaultFuncFuncStmt 会互相覆盖。
// 如果设置了同一个接口的默认数据，会以最后调用为准。
type MockDefaultDataStmt struct {
	m *MockStmt
}

// MockDefaultFuncFuncStmt 用于向 MockStmt 添加指定接口的默认测试方法实现。
// 如果该接口没有设置测试数据或者测试数据已经消费完，这个默认的测试方法就会被调用。
// MockDefaultDataStmt 和 MockDefaultFuncFuncStmt 会互相覆盖。
// 如果设置了同一个接口的默认数据，会以最后调用为准。
type MockDefaultFuncStmt struct {
	m *MockStmt
}

// MockStmt 是一个 mock 容器，可以将 Stmt 接口的所有方法设置 mock 数据。
type MockStmt struct {
	mu sync.Mutex

	data            mockDataStmt
	defData         mockDefaultDataStmt
	impl            mysql.Stmt
	optionalMethods map[string]bool

	dataAdder   *MockDataAdderStmt
	funcAdder   *MockFuncAdderStmt
	defaultData *MockDefaultDataStmt
	defaultFunc *MockDefaultFuncStmt
}

type mockImplStmt struct {
	mysql.Stmt // 强行的实现 Stmt 接口，用于做向后兼容。

	m *MockStmt
}

// 确保 mockImplStmt 始终实现了 mysql.Stmt 接口。
var _ mysql.Stmt = new(mockImplStmt)

// NewStmt 创建一个 Stmt 的 mock 容器。
// impl 是一个默认的接口实现，假如接口没有被 mock 或者 mock 数据已经消费完，impl 的方法会被调用。
//
// impl 可以为 nil，但是这意味着没有默认实现，一旦调用一个没有 mock 数据的接口，且这个接口没有被设置为可选，那么会触发 panic。
func NewStmt(impl mysql.Stmt) *MockStmt {
	optionalMethods := map[string]bool{}
	optionalMethods["Close"] = true
	optionalMethods["Init"] = true

	m := &MockStmt{
		data:            mockDataStmt{},
		defData:         mockDefaultDataStmt{},
		impl:            impl,
		optionalMethods: optionalMethods,
	}
	m.dataAdder = &MockDataAdderStmt{m}
	m.funcAdder = &MockFuncAdderStmt{m}
	m.defaultData = &MockDefaultDataStmt{m}
	m.defaultFunc = &MockDefaultFuncStmt{m}
	return m
}

// Mock 返回一个实现了所有 Stmt 方法的实例。
func (m *MockStmt) Mock() mysql.Stmt {
	return &mockImplStmt{
		m: m,
	}
}

// DataAdder 返回一个添加接口测试数据的入口。
func (m *MockStmt) DataAdder() *MockDataAdderStmt {
	return m.dataAdder
}

// FuncAdder 返回一个添加接口回调的入口。
func (m *MockStmt) FuncAdder() *MockFuncAdderStmt {
	return m.funcAdder
}

// DefaultData 返回一个设置接口默认测试数据的入口。
func (m *MockStmt) DefaultData() *MockDefaultDataStmt {
	return m.defaultData
}

// DefaultFunc 返回一个设置接口默认回调的入口。
func (m *MockStmt) DefaultFunc() *MockDefaultFuncStmt {
	return m.defaultFunc
}

// Close 增加一个 Close 方法的返回值的 mock 数据。
// 如果 Stmt#Close 被调用，这个 mock 数据就会按照 FIFO 顺序被消费掉。
func (__adder *MockDataAdderStmt) Close(err error) {
	__adder.m.mu.Lock()
	defer __adder.m.mu.Unlock()

	__key := "Close"
	__adder.m.data[__key] = append(__adder.m.data[__key], func() error {
		return err
	})
}

// Close 增加一个 mock 数据的回调函数，在执行 Stmt#Close 方法时会调用，并用这个回调函数的返回值来作为接口返回值。
// 如果 Stmt#Close 被调用，这个 mock 数据就会按照 FIFO 顺序被消费掉。
func (adder *MockFuncAdderStmt) Close(f func() (err error)) {
	adder.m.mu.Lock()
	defer adder.m.mu.Unlock()

	key := "Close"
	adder.m.data[key] = append(adder.m.data[key], f)
}

// Close 设置 Stmt#Close 方法的默认 mock 数据，仅当所有 mock 数据消耗完之后起作用。
// 这个函数会覆盖 m.DefaultFunc().Close(...) 方法设置的默认回调。
func (__def *MockDefaultDataStmt) Close(err error) {
	__def.m.mu.Lock()
	defer __def.m.mu.Unlock()

	__key := "Close"
	__def.m.defData[__key] = func() error {
		return err
	}
}

// Close 设置 Stmt#Close 方法的默认回调函数，仅当所有 mock 数据消耗完之后起作用。
// 这个函数会覆盖 m.DefaultData().Close(...) 方法设置的默认数据。
func (def *MockDefaultFuncStmt) Close(f func() (err error)) {
	def.m.mu.Lock()
	defer def.m.mu.Unlock()

	key := "Close"
	def.m.defData[key] = f
}

func (__impl *mockImplStmt) Close() (err error) {
	__impl.m.mu.Lock()

	__key := "Close"
	__data := __impl.m.data[__key]

	if len(__data) == 0 {
		__defData, __ok := __impl.m.defData[__key]
		__impl.m.mu.Unlock()

		if __ok {
			__f := __defData.(func() error)
			return __f()
		}

		if __impl.m.impl != nil {
			return __impl.m.impl.Close()
		}

		if _, __ok := __impl.m.optionalMethods[__key]; __ok {
			return
		}

		panic(fmt.Sprintf("no mock data nor default implementation. [method:%v]", __key))
	}

	// FIFO 顺序。
	__f := __data[0].(func() error)
	__impl.m.data[__key] = __data[1:]
	__impl.m.mu.Unlock()

	return __f()
}

// Exec 增加一个 Exec 方法的返回值的 mock 数据。
// 如果 Stmt#Exec 被调用，这个 mock 数据就会按照 FIFO 顺序被消费掉。
func (__adder *MockDataAdderStmt) Exec(out0 mysql.Result, err error) {
	__adder.m.mu.Lock()
	defer __adder.m.mu.Unlock()

	__key := "Exec"
	__adder.m.data[__key] = append(__adder.m.data[__key], func(args ...interface{}) (mysql.Result, error) {
		return out0, err
	})
}

// Exec 增加一个 mock 数据的回调函数，在执行 Stmt#Exec 方法时会调用，并用这个回调函数的返回值来作为接口返回值。
// 如果 Stmt#Exec 被调用，这个 mock 数据就会按照 FIFO 顺序被消费掉。
func (adder *MockFuncAdderStmt) Exec(f func(args ...interface{}) (out0 mysql.Result, err error)) {
	adder.m.mu.Lock()
	defer adder.m.mu.Unlock()

	key := "Exec"
	adder.m.data[key] = append(adder.m.data[key], f)
}

// Exec 设置 Stmt#Exec 方法的默认 mock 数据，仅当所有 mock 数据消耗完之后起作用。
// 这个函数会覆盖 m.DefaultFunc().Exec(...) 方法设置的默认回调。
func (__def *MockDefaultDataStmt) Exec(out0 mysql.Result, err error) {
	__def.m.mu.Lock()
	defer __def.m.mu.Unlock()

	__key := "Exec"
	__def.m.defData[__key] = func(args ...interface{}) (mysql.Result, error) {
		return out0, err
	}
}

// Exec 设置 Stmt#Exec 方法的默认回调函数，仅当所有 mock 数据消耗完之后起作用。
// 这个函数会覆盖 m.DefaultData().Exec(...) 方法设置的默认数据。
func (def *MockDefaultFuncStmt) Exec(f func(args ...interface{}) (out0 mysql.Result, err error)) {
	def.m.mu.Lock()
	defer def.m.mu.Unlock()

	key := "Exec"
	def.m.defData[key] = f
}

func (__impl *mockImplStmt) Exec(args ...interface{}) (out0 mysql.Result, err error) {
	__impl.m.mu.Lock()

	__key := "Exec"
	__data := __impl.m.data[__key]

	if len(__data) == 0 {
		__defData, __ok := __impl.m.defData[__key]
		__impl.m.mu.Unlock()

		if __ok {
			__f := __defData.(func(args ...interface{}) (mysql.Result, error))
			return __f(args...)
		}

		if __impl.m.impl != nil {
			return __impl.m.impl.Exec(args...)
		}

		if _, __ok := __impl.m.optionalMethods[__key]; __ok {
			return
		}

		panic(fmt.Sprintf("no mock data nor default implementation. [method:%v]", __key))
	}

	// FIFO 顺序。
	__f := __data[0].(func(args ...interface{}) (mysql.Result, error))
	__impl.m.data[__key] = __data[1:]
	__impl.m.mu.Unlock()

	return __f(args...)
}

// Query 增加一个 Query 方法的返回值的 mock 数据。
// 如果 Stmt#Query 被调用，这个 mock 数据就会按照 FIFO 顺序被消费掉。
func (__adder *MockDataAdderStmt) Query(out0 mysql.Rows, err error) {
	__adder.m.mu.Lock()
	defer __adder.m.mu.Unlock()

	__key := "Query"
	__adder.m.data[__key] = append(__adder.m.data[__key], func(args ...interface{}) (mysql.Rows, error) {
		return out0, err
	})
}

// Query 增加一个 mock 数据的回调函数，在执行 Stmt#Query 方法时会调用，并用这个回调函数的返回值来作为接口返回值。
// 如果 Stmt#Query 被调用，这个 mock 数据就会按照 FIFO 顺序被消费掉。
func (adder *MockFuncAdderStmt) Query(f func(args ...interface{}) (out0 mysql.Rows, err error)) {
	adder.m.mu.Lock()
	defer adder.m.mu.Unlock()

	key := "Query"
	adder.m.data[key] = append(adder.m.data[key], f)
}

// Query 设置 Stmt#Query 方法的默认 mock 数据，仅当所有 mock 数据消耗完之后起作用。
// 这个函数会覆盖 m.DefaultFunc().Query(...) 方法设置的默认回调。
func (__def *MockDefaultDataStmt) Query(out0 mysql.Rows, err error) {
	__def.m.mu.Lock()
	defer __def.m.mu.Unlock()

	__key := "Query"
	__def.m.defData[__key] = func(args ...interface{}) (mysql.Rows, error) {
		return out0, err
	}
}

// Query 设置 Stmt#Query 方法的默认回调函数，仅当所有 mock 数据消耗完之后起作用。
// 这个函数会覆盖 m.DefaultData().Query(...) 方法设置的默认数据。
func (def *MockDefaultFuncStmt) Query(f func(args ...interface{}) (out0 mysql.Rows, err error)) {
	def.m.mu.Lock()
	defer def.m.mu.Unlock()

	key := "Query"
	def.m.defData[key] = f
}

func (__impl *mockImplStmt) Query(args ...interface{}) (out0 mysql.Rows, err error) {
	__impl.m.mu.Lock()

	__key := "Query"
	__data := __impl.m.data[__key]

	if len(__data) == 0 {
		__defData, __ok := __impl.m.defData[__key]
		__impl.m.mu.Unlock()

		if __ok {
			__f := __defData.(func(args ...interface{}) (mysql.Rows, error))
			return __f(args...)
		}

		if __impl.m.impl != nil {
			return __impl.m.impl.Query(args...)
		}

		if _, __ok := __impl.m.optionalMethods[__key]; __ok {
			return
		}

		panic(fmt.Sprintf("no mock data nor default implementation. [method:%v]", __key))
	}

	// FIFO 顺序。
	__f := __data[0].(func(args ...interface{}) (mysql.Rows, error))
	__impl.m.data[__key] = __data[1:]
	__impl.m.mu.Unlock()

	return __f(args...)
}

// QueryRow 增加一个 QueryRow 方法的返回值的 mock 数据。
// 如果 Stmt#QueryRow 被调用，这个 mock 数据就会按照 FIFO 顺序被消费掉。
func (__adder *MockDataAdderStmt) QueryRow(out0 mysql.Row, err error) {
	__adder.m.mu.Lock()
	defer __adder.m.mu.Unlock()

	__key := "QueryRow"
	__adder.m.data[__key] = append(__adder.m.data[__key], func(args ...interface{}) (mysql.Row, error) {
		return out0, err
	})
}

// QueryRow 增加一个 mock 数据的回调函数，在执行 Stmt#QueryRow 方法时会调用，并用这个回调函数的返回值来作为接口返回值。
// 如果 Stmt#QueryRow 被调用，这个 mock 数据就会按照 FIFO 顺序被消费掉。
func (adder *MockFuncAdderStmt) QueryRow(f func(args ...interface{}) (out0 mysql.Row, err error)) {
	adder.m.mu.Lock()
	defer adder.m.mu.Unlock()

	key := "QueryRow"
	adder.m.data[key] = append(adder.m.data[key], f)
}

// QueryRow 设置 Stmt#QueryRow 方法的默认 mock 数据，仅当所有 mock 数据消耗完之后起作用。
// 这个函数会覆盖 m.DefaultFunc().QueryRow(...) 方法设置的默认回调。
func (__def *MockDefaultDataStmt) QueryRow(out0 mysql.Row, err error) {
	__def.m.mu.Lock()
	defer __def.m.mu.Unlock()

	__key := "QueryRow"
	__def.m.defData[__key] = func(args ...interface{}) (mysql.Row, error) {
		return out0, err
	}
}

// QueryRow 设置 Stmt#QueryRow 方法的默认回调函数，仅当所有 mock 数据消耗完之后起作用。
// 这个函数会覆盖 m.DefaultData().QueryRow(...) 方法设置的默认数据。
func (def *MockDefaultFuncStmt) QueryRow(f func(args ...interface{}) (out0 mysql.Row, err error)) {
	def.m.mu.Lock()
	defer def.m.mu.Unlock()

	key := "QueryRow"
	def.m.defData[key] = f
}

func (__impl *mockImplStmt) QueryRow(args ...interface{}) (out0 mysql.Row, err error) {
	__impl.m.mu.Lock()

	__key := "QueryRow"
	__data := __impl.m.data[__key]

	if len(__data) == 0 {
		__defData, __ok := __impl.m.defData[__key]
		__impl.m.mu.Unlock()

		if __ok {
			__f := __defData.(func(args ...interface{}) (mysql.Row, error))
			return __f(args...)
		}

		if __impl.m.impl != nil {
			return __impl.m.impl.QueryRow(args...)
		}

		if _, __ok := __impl.m.optionalMethods[__key]; __ok {
			return
		}

		panic(fmt.Sprintf("no mock data nor default implementation. [method:%v]", __key))
	}

	// FIFO 顺序。
	__f := __data[0].(func(args ...interface{}) (mysql.Row, error))
	__impl.m.data[__key] = __data[1:]
	__impl.m.mu.Unlock()

	return __f(args...)
}

// SQLStmt 增加一个 SQLStmt 方法的返回值的 mock 数据。
// 如果 Stmt#SQLStmt 被调用，这个 mock 数据就会按照 FIFO 顺序被消费掉。
func (__adder *MockDataAdderStmt) SQLStmt(out0 *sql.Stmt) {
	__adder.m.mu.Lock()
	defer __adder.m.mu.Unlock()

	__key := "SQLStmt"
	__adder.m.data[__key] = append(__adder.m.data[__key], func() *sql.Stmt {
		return out0
	})
}

// SQLStmt 增加一个 mock 数据的回调函数，在执行 Stmt#SQLStmt 方法时会调用，并用这个回调函数的返回值来作为接口返回值。
// 如果 Stmt#SQLStmt 被调用，这个 mock 数据就会按照 FIFO 顺序被消费掉。
func (adder *MockFuncAdderStmt) SQLStmt(f func() (out0 *sql.Stmt)) {
	adder.m.mu.Lock()
	defer adder.m.mu.Unlock()

	key := "SQLStmt"
	adder.m.data[key] = append(adder.m.data[key], f)
}

// SQLStmt 设置 Stmt#SQLStmt 方法的默认 mock 数据，仅当所有 mock 数据消耗完之后起作用。
// 这个函数会覆盖 m.DefaultFunc().SQLStmt(...) 方法设置的默认回调。
func (__def *MockDefaultDataStmt) SQLStmt(out0 *sql.Stmt) {
	__def.m.mu.Lock()
	defer __def.m.mu.Unlock()

	__key := "SQLStmt"
	__def.m.defData[__key] = func() *sql.Stmt {
		return out0
	}
}

// SQLStmt 设置 Stmt#SQLStmt 方法的默认回调函数，仅当所有 mock 数据消耗完之后起作用。
// 这个函数会覆盖 m.DefaultData().SQLStmt(...) 方法设置的默认数据。
func (def *MockDefaultFuncStmt) SQLStmt(f func() (out0 *sql.Stmt)) {
	def.m.mu.Lock()
	defer def.m.mu.Unlock()

	key := "SQLStmt"
	def.m.defData[key] = f
}

func (__impl *mockImplStmt) SQLStmt() (out0 *sql.Stmt) {
	__impl.m.mu.Lock()

	__key := "SQLStmt"
	__data := __impl.m.data[__key]

	if len(__data) == 0 {
		__defData, __ok := __impl.m.defData[__key]
		__impl.m.mu.Unlock()

		if __ok {
			__f := __defData.(func() *sql.Stmt)
			return __f()
		}

		if __impl.m.impl != nil {
			return __impl.m.impl.SQLStmt()
		}

		if _, __ok := __impl.m.optionalMethods[__key]; __ok {
			return
		}

		panic(fmt.Sprintf("no mock data nor default implementation. [method:%v]", __key))
	}

	// FIFO 顺序。
	__f := __data[0].(func() *sql.Stmt)
	__impl.m.data[__key] = __data[1:]
	__impl.m.mu.Unlock()

	return __f()
}
