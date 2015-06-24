package gomodel

import "database/sql"

type Executor interface {
	Exec(...interface{}) (sql.Result, error)
	Query(...interface{}) (*sql.Rows, error)
	Close() error
}

// Update always returl the count of affected rows
func Update(exec Executor, err error, args ...interface{}) (int64, error) {
	return Exec(exec, err, RES_ROWS, args...)
}

// Exec execute stmt with given arguments and resolve the result if error is nil
func Exec(exec Executor, err error, typ ResultType, args ...interface{}) (int64, error) {
	if err != nil {
		return 0, err
	}

	res, err := exec.Exec(args...)
	return ResolveResult(res, err, typ)
}

// Query execute the query stmt, error stored in Scanner
func Query(exec Executor, err error, args ...interface{}) (Scanner, *sql.Rows) {
	if err != nil {
		return Scanner{err}, nil
	}

	rows, err := exec.Query(args...)
	if err != nil {
		return Scanner{err}, nil
	}

	return normalScanner, rows
}

// Update always returl the count of affected rows
func CloseUpdate(exec Executor, err error, args ...interface{}) (int64, error) {
	return CloseExec(exec, err, RES_ROWS, args...)
}

// Exec execute stmt with given arguments and resolve the result if error is nil
func CloseExec(exec Executor, err error, typ ResultType, args ...interface{}) (int64, error) {
	c, err := Exec(exec, err, typ, args...)
	exec.Close()

	return c, err
}

// Query execute the query stmt, error stored in Scanner
func CloseQuery(exec Executor, err error, args ...interface{}) (Scanner, *sql.Rows) {
	sc, rows := Query(exec, err, args...)
	exec.Close()

	return sc, rows
}
